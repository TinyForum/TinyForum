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

// PluginManifest 表示插件压缩包内 manifest.json 的结构
type PluginManifest struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Type        do.PluginType     `json:"type"`
	Category    do.PluginCategory `json:"category"`
	Entry       string            `json:"entry"` // 插件入口文件（相对路径）
}

// Validate 校验 manifest 必填字段
func (m *PluginManifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("%w: missing name", apperrors.ErrInvalidManifest)
	}
	if m.Version == "" {
		return fmt.Errorf("%w: missing version", apperrors.ErrInvalidManifest)
	}
	if m.Entry == "" {
		return fmt.Errorf("%w: missing entry", apperrors.ErrInvalidManifest)
	}
	return nil
}

const (
	defaultStorageDir = "./plugins"
	zipExtension      = ".zip"
	manifestFileName  = "manifest.json"
	macOSMetaPrefix   = "__MACOSX"
	dsStoreFile       = ".DS_Store"
	dirPerm           = 0755
)

// Create 上传并安装插件（若同名插件已存在则覆盖）
//
// 流程：
//  1. 校验文件格式（仅允许 .zip）
//  2. 保存到临时文件
//  3. 解析 manifest.json
//  4. 查询同名插件（含软删除）
//  5. 解压插件到持久化目录
//  6. 写入或更新数据库记录（失败时回滚文件）
func (s *pluginService) Create(ctx context.Context, fileHeader *multipart.FileHeader, userID uint) (*do.PluginMeta, error) {
	// 1. 校验文件扩展名
	if !strings.HasSuffix(fileHeader.Filename, zipExtension) {
		return nil, apperrors.ErrInvalidPluginFormat
	}

	// 2. 保存上传文件到系统临时目录
	tempZipPath, cleanupTemp, err := s.saveToTempFile(fileHeader)
	if err != nil {
		return nil, err
	}
	defer cleanupTemp()

	// 3. 从 ZIP 中解析 manifest.json
	manifest, manifestDir, err := s.parseManifestFromZip(tempZipPath)
	if err != nil {
		return nil, err
	}

	// 4. 查询同名插件（包含软删除记录，用于判断是否覆盖）
	existing, err := s.findExistingPlugin(ctx, manifest.Name)
	if err != nil {
		return nil, err
	}

	// 5. 确定目标解压目录（格式：<storageBase>/<name>/<version>/）
	targetDir := s.buildTargetDir(manifest.Slug, manifest.Version)

	// 6. 解压插件文件到目标目录
	if err := s.extractPluginFiles(tempZipPath, targetDir, manifestDir); err != nil {
		logger.Errorf("解压插件文件失败: %v", err)
		return nil, err
	}

	// 7. 写入或更新数据库记录；失败时回滚已解压的文件
	pluginMeta, err := s.saveOrUpdatePluginMeta(ctx, existing, manifest, targetDir, userID)
	if err != nil {
		if rollbackErr := s.storage.DeleteDir(targetDir); rollbackErr != nil {
			logger.Errorf("回滚插件目录失败 [%s]: %v", targetDir, rollbackErr)
		}
		return nil, err
	}

	return pluginMeta, nil
}

// ── 内部辅助方法 ────────────────────────────────────────────────────────────────

// saveToTempFile 将 multipart 上传文件落盘到系统临时目录，返回临时路径和清理函数。
// 调用方必须 defer 执行清理函数以删除临时文件。
func (s *pluginService) saveToTempFile(fileHeader *multipart.FileHeader) (tmpPath string, cleanup func(), err error) {
	src, err := fileHeader.Open()
	if err != nil {
		return "", nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer src.Close()

	tmpFile, err := os.CreateTemp("", "plugin-*.zip")
	if err != nil {
		return "", nil, fmt.Errorf("创建临时文件失败: %w", err)
	}

	cleanup = func() {
		tmpFile.Close()
		if removeErr := os.Remove(tmpFile.Name()); removeErr != nil && !os.IsNotExist(removeErr) {
			logger.Warnf("清理临时文件失败 [%s]: %v", tmpFile.Name(), removeErr)
		}
	}

	if _, err = io.Copy(tmpFile, src); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("写入临时文件失败: %w", err)
	}

	// 重置读指针供后续 zip.OpenReader 使用
	if _, err = tmpFile.Seek(0, io.SeekStart); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("重置临时文件指针失败: %w", err)
	}

	return tmpFile.Name(), cleanup, nil
}

// parseManifestFromZip 从 ZIP 归档中定位并解析 manifest.json。
// 返回解析结果、manifest 所在目录（相对于 ZIP 根）和错误。
func (s *pluginService) parseManifestFromZip(zipPath string) (PluginManifest, string, error) {
	logger.Infof("解析 ZIP 文件: %s", zipPath)

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return PluginManifest{}, "", fmt.Errorf("打开 ZIP 文件失败: %w", err)
	}
	defer r.Close()

	manifestFile, manifestDir, err := findManifestInZip(r.File)
	if err != nil {
		return PluginManifest{}, "", err
	}

	manifest, err := readAndParseManifest(manifestFile)
	if err != nil {
		return PluginManifest{}, "", err
	}

	if err := manifest.Validate(); err != nil {
		return PluginManifest{}, "", err
	}

	return manifest, manifestDir, nil
}

// findManifestInZip 遍历 ZIP 文件列表，找到 manifest.json 所在的条目和其父目录。
func findManifestInZip(files []*zip.File) (*zip.File, string, error) {
	for _, f := range files {
		if f.FileInfo().IsDir() {
			continue
		}
		cleanName := filepath.Clean(f.Name)
		if isSystemFile(cleanName) {
			continue
		}
		if filepath.Base(cleanName) == manifestFileName {
			dir := filepath.Dir(cleanName)
			if dir == "." {
				dir = ""
			}
			return f, dir, nil
		}
	}
	return nil, "", apperrors.ErrManifestNotFound
}

// readAndParseManifest 从 zip.File 中读取并反序列化 manifest.json。
func readAndParseManifest(f *zip.File) (PluginManifest, error) {
	rc, err := f.Open()
	if err != nil {
		return PluginManifest{}, fmt.Errorf("打开 manifest.json 失败: %w", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return PluginManifest{}, fmt.Errorf("读取 manifest.json 失败: %w", err)
	}

	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return PluginManifest{}, fmt.Errorf("解析 manifest.json 失败: %w", err)
	}

	return manifest, nil
}

// extractPluginFiles 将 ZIP 解压到 targetDir，自动剥离 manifestDir 前缀，
// 并执行 Zip Slip 路径安全检查。
func (s *pluginService) extractPluginFiles(zipPath, targetDir, manifestDir string) error {
	logger.Infof("解压插件到目标路径: %s", targetDir)

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("打开 ZIP 文件失败: %w", err)
	}
	defer r.Close()

	// 清理目标目录，保证幂等（重复安装同一版本时覆盖旧文件）
	if err := s.storage.DeleteDir(targetDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("清理目标目录失败: %w", err)
	}

	for _, f := range r.File {
		if err := s.extractSingleFile(f, targetDir, manifestDir); err != nil {
			return err
		}
	}
	return nil
}

// extractSingleFile 解压 ZIP 中的单个文件条目到 targetDir。
func (s *pluginService) extractSingleFile(f *zip.File, targetDir, manifestDir string) error {
	if f.FileInfo().IsDir() {
		return nil
	}

	cleanName := filepath.Clean(f.Name)

	// 跳过 macOS 系统元数据
	if isSystemFile(cleanName) {
		return nil
	}

	// 安全检查：原始路径不可含路径穿越
	if containsPathTraversal(cleanName) {
		return fmt.Errorf("不合法的文件路径（疑似路径穿越）: %s", f.Name)
	}

	// 剥离 manifestDir 前缀，得到相对于插件根的路径
	relPath := stripManifestPrefix(cleanName, manifestDir)
	if relPath == "" {
		return nil // 跳过 manifest 目录本身
	}
	if containsPathTraversal(relPath) {
		return fmt.Errorf("不合法的相对路径: %s", relPath)
	}

	destPath := filepath.Join(targetDir, relPath)

	// 二次安全检查：最终落盘路径必须在 targetDir 内
	if !strings.HasPrefix(destPath, targetDir+string(os.PathSeparator)) {
		return fmt.Errorf("路径逃逸检测: %s", f.Name)
	}

	// 创建父目录
	if err := os.MkdirAll(filepath.Dir(destPath), dirPerm); err != nil {
		return fmt.Errorf("创建目录失败 [%s]: %w", filepath.Dir(destPath), err)
	}

	// 写入文件
	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("打开 ZIP 条目失败 [%s]: %w", f.Name, err)
	}
	defer rc.Close()

	logger.Infof("写入插件文件: %s", destPath)
	if _, err := s.storage.Save(rc, destPath); err != nil {
		return fmt.Errorf("保存插件文件失败 [%s]: %w", f.Name, err)
	}
	return nil
}

// saveOrUpdatePluginMeta 根据是否存在同名记录，执行插入或更新操作。
func (s *pluginService) saveOrUpdatePluginMeta(
	ctx context.Context,
	existing *do.PluginMeta,
	manifest PluginManifest,
	targetDir string,
	userID uint,
) (*do.PluginMeta, error) {
	logger.Infof("持久化插件元数据: id=%s name=%s version=%s", manifest.Slug, manifest.Name, manifest.Version)

	// ScriptURL 格式：/store/plugins/<id>/<version>/<entry>
	scriptURL := "/store/plugins/" + path.Join(manifest.Slug, manifest.Version, manifest.Entry)

	if existing != nil {
		return s.overwritePluginMeta(ctx, existing, manifest, targetDir, scriptURL, userID)
	}
	return s.createPluginMeta(ctx, manifest, scriptURL, userID)
}

// overwritePluginMeta 更新已有插件记录，并在版本变更时删除旧版本目录。
func (s *pluginService) overwritePluginMeta(
	ctx context.Context,
	existing *do.PluginMeta,
	manifest PluginManifest,
	newTargetDir string,
	scriptURL string,
	userID uint,
) (*do.PluginMeta, error) {
	// 版本不同时，删除旧版本的物理目录
	oldDir := s.buildTargetDir(existing.Name, existing.Version)
	if oldDir != newTargetDir {
		if err := s.storage.DeleteDir(oldDir); err != nil && !os.IsNotExist(err) {
			logger.Warnf("删除旧版本目录失败 [%s]: %v", oldDir, err)
			// 非致命错误，继续执行
		}
	}

	// 更新记录字段（同时恢复软删除状态）
	existing.DeletedAt = gorm.DeletedAt{}
	existing.Version = manifest.Version
	existing.Description = manifest.Description
	existing.Type = manifest.Type
	existing.Category = manifest.Category
	existing.ScriptURL = scriptURL
	existing.Enabled = false
	existing.Status = do.PluginStatusInactive
	existing.AuthorID = userID

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("更新插件元数据失败: %w", err)
	}

	logger.Infof("插件已覆盖安装: name=%s version=%s userID=%d", manifest.Name, manifest.Version, userID)
	return existing, nil
}

// createPluginMeta 创建全新的插件元数据记录。
func (s *pluginService) createPluginMeta(
	ctx context.Context,
	manifest PluginManifest,
	scriptURL string,
	userID uint,
) (*do.PluginMeta, error) {
	pluginMeta := &do.PluginMeta{
		Name:        manifest.Name,
		Slug:        manifest.Slug,
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
		return nil, fmt.Errorf("创建插件元数据失败: %w", err)
	}

	logger.Infof("新插件安装成功: name=%s version=%s userID=%d", manifest.Name, manifest.Version, userID)
	return pluginMeta, nil
}

// findExistingPlugin 查询同名插件（含软删除），未找到时返回 nil, nil。
func (s *pluginService) findExistingPlugin(ctx context.Context, name string) (*do.PluginMeta, error) {
	existing, err := s.repo.FindByNameUnscoped(ctx, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("查询已有插件失败: %w", err)
	}
	return existing, nil
}

// buildTargetDir 拼接插件的持久化存储目录路径。
// 格式：<storageBase>/<name>/<version>
func (s *pluginService) buildTargetDir(name, version string) string {
	return filepath.Join(s.getStorageBase(), name, version)
}

// getStorageBase 返回插件存储根目录，优先使用配置值。
func (s *pluginService) getStorageBase() string {
	if s.cfg != nil && s.cfg.StorageDir != "" {
		return s.cfg.StorageDir
	}
	return defaultStorageDir
}

// ── 纯函数工具 ──────────────────────────────────────────────────────────────────

// isSystemFile 判断文件是否为应跳过的系统元数据（macOS __MACOSX / .DS_Store）。
func isSystemFile(cleanPath string) bool {
	return strings.HasPrefix(cleanPath, macOSMetaPrefix) ||
		strings.HasSuffix(cleanPath, dsStoreFile)
}

// containsPathTraversal 判断路径是否包含路径穿越片段（".."）。
func containsPathTraversal(p string) bool {
	return strings.Contains(p, "..")
}

// stripManifestPrefix 从完整路径中剥离 manifestDir 前缀，返回相对于插件根的路径。
func stripManifestPrefix(cleanName, manifestDir string) string {
	if manifestDir == "" {
		return cleanName
	}
	prefix := manifestDir + string(os.PathSeparator)
	if strings.HasPrefix(cleanName, prefix) {
		return strings.TrimPrefix(cleanName, prefix)
	}
	return cleanName
}
