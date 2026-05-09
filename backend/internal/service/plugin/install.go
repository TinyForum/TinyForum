// // service/plugin/install.go
// package plugin

// import (
// 	"archive/zip"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"path/filepath"
// 	"strings"
// 	"tiny-forum/internal/model/do"
// )

// // Install 实现
// func (s *pluginService) Install(ctx context.Context, fileHeader *multipart.FileHeader, userID uint) (*do.PluginMeta, error) {
//     // 1. 校验文件扩展名
//     if !strings.HasSuffix(fileHeader.Filename, ".zip") {
//         return nil, fmt.Errorf("only zip files are allowed")
//     }

//     // 2. 打开上传的文件流
//     src, err := fileHeader.Open()
//     if err != nil {
//         return nil, fmt.Errorf("open uploaded file: %w", err)
//     }
//     defer src.Close()

//     // 3. 读取 zip 内容到内存（考虑到插件包通常不大，几十MB以内）
//     zipData, err := io.ReadAll(src)
//     if err != nil {
//         return nil, fmt.Errorf("read zip data: %w", err)
//     }
//     zipReader, err := zip.NewReader(strings.NewReader(string(zipData)), int64(len(zipData)))
//     if err != nil {
//         return nil, fmt.Errorf("invalid zip file: %w", err)
//     }

//     // 4. 查找 manifest.json 文件
//     var manifestFile *zip.File
//     for _, f := range zipReader.File {
//         if f.Name == "manifest.json" {
//             manifestFile = f
//             break
//         }
//     }
//     if manifestFile == nil {
//         return nil, fmt.Errorf("manifest.json not found in plugin package")
//     }

//     // 5. 读取并解析 manifest.json
//     rc, err := manifestFile.Open()
//     if err != nil {
//         return nil, fmt.Errorf("open manifest.json: %w", err)
//     }
//     defer rc.Close()
//     manifestData, err := io.ReadAll(rc)
//     if err != nil {
//         return nil, fmt.Errorf("read manifest.json: %w", err)
//     }
//     var manifest struct {
//         Name        string          `json:"name"`
//         Version     string          `json:"version"`
//         Description string          `json:"description"`
//         Type        do.PluginType   `json:"type"`
//         Category    do.PluginCategory `json:"category"`
//         Entry       string          `json:"entry"` // 前端入口文件相对路径
//         // 其他字段可根据需要扩展
//     }
//     if err := json.Unmarshal(manifestData, &manifest); err != nil {
//         return nil, fmt.Errorf("parse manifest.json: %w", err)
//     }

//     // 6. 检查插件是否已存在（根据 name）
//     existing, _ := s.repo.GetByName(ctx, manifest.Name)
//     if existing != nil {
//         return nil, fmt.Errorf("plugin %s already exists", manifest.Name)
//     }

//     // 7. 生成插件存储目录（基于 name + version）
//     pluginDir := filepath.Join("plugins", manifest.Name, manifest.Version)

//     // 8. 将 zip 中的所有文件解压到持久化存储（通过 StorageDriver）
//     for _, f := range zipReader.File {
//         if f.FileInfo().IsDir() {
//             continue // 目录无需单独存储，由文件上级路径自动创建
//         }
//         destPath := filepath.Join(pluginDir, f.Name)
//         rc, err := f.Open()
//         if err != nil {
//             return nil, fmt.Errorf("open file in zip: %w", err)
//         }
//         // 使用 StorageDriver 保存文件
//         if _, err := s.storage.Save(rc, destPath); err != nil {
//             rc.Close()
//             return nil, fmt.Errorf("save plugin file %s: %w", f.Name, err)
//         }
//         rc.Close()
//     }

//     // 9. 构建 PluginMeta 记录
//     pluginMeta := &do.PluginMeta{
//         Name:        manifest.Name,
//         Version:     manifest.Version,
//         Description: manifest.Description,
//         Type:        manifest.Type,
//         Category:    manifest.Category,
//         ScriptURL:   "/plugin-files/" + filepath.Join(pluginDir, manifest.Entry),
//         AuthorID:    userID,
//         Enabled:     false,   // 安装后默认未启用，需要管理员或用户手动启用
//         Status:      do.PluginStatusInactive,
//         // 其他字段使用默认值
//     }

//     // 10. 写入数据库
//     if err := s.repo.Create(ctx, pluginMeta); err != nil {
//         // 创建失败，清理已存储的文件
//         _ = s.storage.DeleteDir(pluginDir) // 需要 StorageDriver 支持递归删除
//         return nil, fmt.Errorf("save plugin metadata: %w", err)
//     }

//	    return pluginMeta, nil
//	}
package plugin