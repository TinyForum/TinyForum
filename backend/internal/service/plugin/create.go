package plugin

import (
	"archive/zip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"gorm.io/gorm"

	"tiny-forum/internal/model/do"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
)

type PluginManifest struct {
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Type        do.PluginType     `json:"type"`
    Category    do.PluginCategory `json:"category"`
    Entry       string            `json:"entry"`
}
// Create 安装或覆盖插件
func (s *pluginService) Create(ctx context.Context, fileHeader *multipart.FileHeader, userID uint) (*do.PluginMeta, error) {
    // 1. 校验文件扩展名
    if !strings.HasSuffix(fileHeader.Filename, ".zip") {
        return nil, apperrors.ErrInvalidPluginFormat
    }

    // 2. 保存上传文件到临时 ZIP 文件
    tempZipPath, cleanup, err := s.saveToTempFile(fileHeader)
    if err != nil {
        return nil, err
    }
    defer cleanup()

    // 3. 解析 ZIP 并读取 manifest
    manifest, manifestDir, err := s.parseManifestFromZip(tempZipPath)
    if err != nil {
        return nil, err
    }

    // 4. 检查插件是否已存在（包括软删除）
    existing, err := s.repo.FindByNameUnscoped(ctx, manifest.Name)
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, fmt.Errorf("check existing plugin: %w", err)
    }

    // 5. 准备目标存储目录（基于 name 和 version）
    storageBase := s.getStorageBase()
    targetDir := filepath.Join(storageBase, manifest.Name, manifest.Version)

    // 6. 解压插件文件
    if err := s.extractPluginFiles(tempZipPath, targetDir, manifestDir); err != nil {
        return nil, err
    }

    // 7. 存储或更新数据库记录
    pluginMeta, err := s.saveOrUpdatePluginMeta(ctx, existing, manifest, targetDir, userID)
    if err != nil {
        // 回滚已解压的文件
        _ = s.storage.DeleteDir(targetDir)
        return nil, err
    }

    return pluginMeta, nil
}

// saveToTempFile 将上传的文件保存到临时 ZIP 文件，返回路径和清理函数
func (s *pluginService) saveToTempFile(fileHeader *multipart.FileHeader) (string, func(), error) {
    src, err := fileHeader.Open()
    if err != nil {
        return "", nil, fmt.Errorf("open uploaded file: %w", err)
    }
    defer src.Close()

    tempFile, err := os.CreateTemp("", "plugin-*.zip")
    if err != nil {
        return "", nil, fmt.Errorf("create temp file: %w", err)
    }

    cleanup := func() {
        tempFile.Close()
        os.Remove(tempFile.Name())
    }

    if _, err := io.Copy(tempFile, src); err != nil {
        cleanup()
        return "", nil, fmt.Errorf("copy to temp file: %w", err)
    }

    // 重置文件指针以便后续读取
    if _, err := tempFile.Seek(0, 0); err != nil {
        cleanup()
        return "", nil, fmt.Errorf("seek temp file: %w", err)
    }

    return tempFile.Name(), cleanup, nil
}

// parseManifestFromZip 从 ZIP 文件中读取并解析 manifest.json
func (s *pluginService) parseManifestFromZip(zipPath string) (PluginManifest, string, error) {
    r, err := zip.OpenReader(zipPath)
    if err != nil {
        return PluginManifest{}, "", fmt.Errorf("open zip: %w", err)
    }
    defer r.Close()

    type manifestData struct {
        Name        string            `json:"name"`
        Version     string            `json:"version"`
        Description string            `json:"description"`
        Type        do.PluginType     `json:"type"`
        Category    do.PluginCategory `json:"category"`
        Entry       string            `json:"entry"`
    }

    var manifest PluginManifest
    var manifestFile *zip.File
    var manifestDir string

    for _, f := range r.File {
        if f.FileInfo().IsDir() {
            continue
        }
        cleanName := filepath.Clean(f.Name)
        // 跳过 macOS 元数据
        if strings.HasPrefix(cleanName, "__MACOSX") || strings.HasSuffix(cleanName, ".DS_Store") {
            continue
        }
        if filepath.Base(cleanName) == "manifest.json" {
            manifestFile = f
            manifestDir = filepath.Dir(cleanName)
            if manifestDir == "." {
                manifestDir = ""
            }
            break
        }
    }

    if manifestFile == nil {
        return PluginManifest{}, "", apperrors.ErrManifestNotFound
    }

    rc, err := manifestFile.Open()
    if err != nil {
        return PluginManifest{}, "", fmt.Errorf("open manifest.json: %w", err)
    }
    defer rc.Close()

    data, err := io.ReadAll(rc)
    if err != nil {
        return PluginManifest{}, "", fmt.Errorf("read manifest.json: %w", err)
    }

    if err := json.Unmarshal(data, &manifest); err != nil {
        return PluginManifest{}, "", fmt.Errorf("parse manifest.json: %w", err)
    }

    if manifest.Name == "" || manifest.Version == "" || manifest.Entry == "" {
        return PluginManifest{}, "", apperrors.ErrInvalidManifest
    }

    return manifest, manifestDir, nil
}

// extractPluginFiles 将 ZIP 解压到目标目录，自动去除 manifestDir 前缀
func (s *pluginService) extractPluginFiles(zipPath, targetDir, manifestDir string) error {
    r, err := zip.OpenReader(zipPath)
    if err != nil {
        return fmt.Errorf("open zip for extract: %w", err)
    }
    defer r.Close()

    // 清理目标目录（确保干净解压）
    if err := s.storage.DeleteDir(targetDir); err != nil && !os.IsNotExist(err) {
        return fmt.Errorf("clean target dir: %w", err)
    }

    for _, f := range r.File {
        if f.FileInfo().IsDir() {
            continue
        }
        cleanName := filepath.Clean(f.Name)
        if strings.HasPrefix(cleanName, "__MACOSX") || strings.HasSuffix(cleanName, ".DS_Store") {
            continue
        }
        // 防止 Zip Slip
        if strings.Contains(cleanName, "..") {
            return fmt.Errorf("invalid file path: %s", f.Name)
        }

        // 相对路径：去除 manifestDir 前缀
        relPath := cleanName
        if manifestDir != "" && strings.HasPrefix(cleanName, manifestDir+string(os.PathSeparator)) {
            relPath = strings.TrimPrefix(cleanName, manifestDir+string(os.PathSeparator))
        }
        if relPath == "" {
            continue
        }
        if strings.Contains(relPath, "..") {
            return fmt.Errorf("invalid relative path: %s", relPath)
        }

        destPath := filepath.Join(targetDir, relPath)
        if !strings.HasPrefix(destPath, targetDir+string(os.PathSeparator)) {
            return fmt.Errorf("path escape: %s", f.Name)
        }

        // 创建父目录
        if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
            return fmt.Errorf("create parent dir: %w", err)
        }

        rc, err := f.Open()
        if err != nil {
            return fmt.Errorf("open zip file %s: %w", f.Name, err)
        }
        if _, err := s.storage.Save(rc, destPath); err != nil {
            rc.Close()
            return fmt.Errorf("save file %s: %w", f.Name, err)
        }
        rc.Close()
    }
    return nil
}

// saveOrUpdatePluginMeta 根据现有记录决定插入或更新数据库
func (s *pluginService) saveOrUpdatePluginMeta(ctx context.Context, existing *do.PluginMeta, manifest PluginManifest, targetDir string, userID uint) (*do.PluginMeta, error) {
    scriptURL := "/plugins/" + path.Join(manifest.Name, manifest.Version, manifest.Entry)

    if existing != nil {
        // 覆盖：删除旧物理文件（如果版本不同，旧目录会被新目录覆盖；版本相同时，文件已覆盖）
        oldDir := filepath.Join(s.getStorageBase(), existing.Name, existing.Version)
        if oldDir != targetDir {
            // 不同版本，删除旧目录
            _ = s.storage.DeleteDir(oldDir)
        }

        // 更新所有可变字段
        existing.DeletedAt = gorm.DeletedAt{}        // 恢复软删除（如果有）
        existing.Version = manifest.Version
        existing.Description = manifest.Description
        existing.Type = manifest.Type
        existing.Category = manifest.Category
        existing.ScriptURL = scriptURL
        existing.Enabled = false
        existing.Status = do.PluginStatusInactive
        // 可选：是否更新作者？保持原作者或设为当前用户，根据业务决定
        // 此处保持原有作者，如需更改则取消注释下一行
        // existing.AuthorID = userID

        if err := s.repo.Update(ctx, existing); err != nil {
            return nil, fmt.Errorf("update plugin meta: %w", err)
        }
        logger.Infof("Plugin %s version %s overwritten by user %d", manifest.Name, manifest.Version, userID)
        return existing, nil
    }

    // 新建
    pluginMeta := &do.PluginMeta{
        Name:        manifest.Name,
        Version:     manifest.Version,
        Description: manifest.Description,
        Type:        manifest.Type,
        Category:    manifest.Category,
        ScriptURL:   scriptURL,
        Enabled:     false,
        Status:      do.PluginStatusInactive,
        AuthorID:    userID,
    }
    if err := s.repo.Create(ctx, pluginMeta); err != nil {
        return nil, fmt.Errorf("save plugin meta: %w", err)
    }
    logger.Infof("New plugin %s version %s installed by user %d", manifest.Name, manifest.Version, userID)
    return pluginMeta, nil
}

// getStorageBase 返回插件存储根目录（带默认值）
func (s *pluginService) getStorageBase() string {
    if s.cfg != nil && s.cfg.StorageDir != "" {
        return s.cfg.StorageDir
    }
    return "./plugins" // 默认路径
}