package plugin

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
	pluginRepo "tiny-forum/internal/repository/plugin"
	"tiny-forum/internal/storage"

	"github.com/google/uuid"
)

type pluginService struct {
	repo    pluginRepo.PluginRepository
	storage storage.StorageDriver
}
type PluginService interface {
   ListPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[vo.PluginMetaVO], error) 
   Install(ctx context.Context, fileHeader *multipart.FileHeader, userID uint) (*do.PluginMeta, error)
    ListUserPlugins(ctx context.Context, userID int64) ([]*do.PluginMeta, error)
}

func NewPluginService(repo pluginRepo.PluginRepository, storage storage.StorageDriver) PluginService {
	return &pluginService{repo: repo, storage: storage}
}

// ListPlugins 分页查询插件（实现略，可参考之前）
func (s *pluginService) ListPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[vo.PluginMetaVO], error) {
	// 此处简单返回空，实际需调用 repo 分页
	return nil, nil
}

// Install 安装插件（完整实现）
func (s *pluginService) Install(ctx context.Context, fileHeader *multipart.FileHeader, userID uint) (*do.PluginMeta, error) {
	// 1. 校验扩展名
	if !strings.HasSuffix(fileHeader.Filename, ".zip") {
		return nil, fmt.Errorf("only zip files are allowed")
	}

	// 2. 保存临时 zip 文件
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	tempZipName := uuid.New().String() + ".zip"
	tempZipPath := filepath.Join("temp", tempZipName)
	if _, err := s.storage.Save(src, tempZipPath); err != nil {
		return nil, fmt.Errorf("save temp zip: %w", err)
	}
	defer s.storage.Delete(tempZipPath)

	// 3. 读取 zip 内容到内存（或者先下载到本地临时文件，出于演示使用内存）
	// 重新打开以读取内容
	src2, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src2.Close()

	// 获取文件大小
	stat, _ := fileHeader.Open()
		// 获取文件大小
	size := fileHeader.Size

	stat.Close()

	zipReader, err := zip.NewReader(src2, size)
	if err != nil {
		return nil, fmt.Errorf("read zip: %w", err)
	}

	// 4. 解析 manifest.json
	var manifest struct {
		Name        string          `json:"name"`
		Version     string          `json:"version"`
		Description string          `json:"description"`
		Type        do.PluginType   `json:"type"`
		Category    do.PluginCategory `json:"category"`
		Entry       string          `json:"entry"` // 前端入口文件
	}
	found := false
	for _, f := range zipReader.File {
		if f.Name == "manifest.json" {
			rc, err := f.Open()
			if err != nil {
				return nil, err
			}
			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(data, &manifest); err != nil {
				return nil, fmt.Errorf("parse manifest: %w", err)
			}
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("manifest.json not found in zip")
	}
	if manifest.Name == "" || manifest.Version == "" {
		return nil, fmt.Errorf("invalid manifest: name and version required")
	}

	// 5. 检查插件是否已存在
	existing, _ := s.repo.GetByName(ctx, manifest.Name)
	if existing != nil {
		return nil, fmt.Errorf("plugin %s already exists", manifest.Name)
	}

	// 6. 将插件文件解压到持久化目录
	pluginStoragePrefix := filepath.Join("plugins", manifest.Name, manifest.Version)
	// 先删除旧目录（如果存在）
	_ = s.storage.DeleteDir(pluginStoragePrefix)

	// 遍历 zip 文件并保存到 storage
	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		destPath := filepath.Join(pluginStoragePrefix, f.Name)
		if _, err := s.storage.Save(rc, destPath); err != nil {
			rc.Close()
			// 清理已保存的文件
			s.storage.DeleteDir(pluginStoragePrefix)
			return nil, fmt.Errorf("save plugin file %s: %w", f.Name, err)
		}
		rc.Close()
	}

	// 7. 写入数据库
	pluginMeta := &do.PluginMeta{
		Name:        manifest.Name,
		Version:     manifest.Version,
		Description: manifest.Description,
		Type:        manifest.Type,
		Category:    manifest.Category,
		ScriptURL:   "/plugin-files/" + filepath.Join(pluginStoragePrefix, manifest.Entry),
		Enabled:     false,
		Status:      do.PluginStatusInactive,
		AuthorID:    userID,
	
	}
	if err := s.repo.Create(ctx, pluginMeta); err != nil {
		// 清理已存储的文件
		s.storage.DeleteDir(pluginStoragePrefix)
		return nil, fmt.Errorf("save plugin meta: %w", err)
	}
	return pluginMeta, nil
}

// ListUserPlugins 获取用户安装的插件
func (s *pluginService) ListUserPlugins(ctx context.Context, userID int64) ([]*do.PluginMeta, error) {
	return s.repo.ListByAuthorID(ctx, userID)
}