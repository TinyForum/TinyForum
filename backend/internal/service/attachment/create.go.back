// package attachment

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"github.com/disintegration/imaging"

// 	"tiny-forum/internal/model/bo"
// 	"tiny-forum/internal/model/do"
// 	"tiny-forum/internal/model/dto"
// 	"tiny-forum/internal/model/request"
// )

// // UploadPlugin 专门处理插件上传
// func (s *service) UploadPlugin(ctx context.Context, request bo.PluginUpdateBO) (*dto.UploadResponse, error) {
// 	fileHeader := request.FileHeader
// 	userID := request.UserID

// 	// 1. 验证插件文件（仅 ZIP）
// 	ext, err := s.validatePluginFile(fileHeader)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 2. 生成存储路径
// 	fileID := generateFileID()
// 	storedName := fmt.Sprintf("%s%s", fileID, ext) // 例如：abc123.zip
// 	saveDir := filepath.Join(s.uploadDir, pluginUploadSubDir)
// 	if err := os.MkdirAll(saveDir, 0755); err != nil {
// 		return nil, fmt.Errorf("创建插件目录失败: %w", err)
// 	}
// 	storedPath := filepath.Join(saveDir, storedName)

// 	// 3. 保存文件
// 	srcFile, err := fileHeader.Open()
// 	if err != nil {
// 		return nil, fmt.Errorf("打开插件文件失败: %w", err)
// 	}
// 	defer srcFile.Close()

// 	dstFile, err := os.Create(storedPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("创建插件文件失败: %w", err)
// 	}
// 	defer dstFile.Close()

// 	if _, err := io.Copy(dstFile, srcFile); err != nil {
// 		os.Remove(storedPath)
// 		return nil, fmt.Errorf("保存插件文件失败: %w", err)
// 	}

// 	// 4. 校验 ZIP 内容

// 	// 检查 manifest.json 	格式

// 	// 5. 保存到附件表（标记 file_type = plugin）
// 	attachment := &do.Attachment{
// 		FileID:       fileID,
// 		UserID:       int64(userID),
// 		OriginalName: fileHeader.Filename,
// 		StoredName:   storedName,
// 		StoredPath:   storedPath,
// 		Size:         fileHeader.Size,
// 		MimeType:     "application/zip",
// 		FileType:     "plugin", // 明确标记为插件
// 		Ext:          ext,
// 		Status:       1,
// 		// PostID/CommentID 等留空
// 	}
// 	if err := s.repo.Create(ctx, attachment); err != nil {
// 		os.Remove(storedPath)
// 		return nil, fmt.Errorf("保存插件记录失败: %w", err)
// 	}

// 	// 6. 构建响应 URL
// 	url := fmt.Sprintf("%s/%s/%s", s.urlPrefix, pluginUploadSubDir, storedName)
// 	return &dto.UploadResponse{
// 		FileID:       fileID,
// 		URL:          url,
// 		OriginalName: fileHeader.Filename,
// 		Size:         fileHeader.Size,
// 		MimeType:     "application/zip",
// 	}, nil
// }

// // 主上传方法
// func (s *service) UploadFile(ctx context.Context, userID int64, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest) (*dto.UploadResponse, error) {
// 	// 1. 验证文件
// 	ext, err := s.validateFile(fileHeader, req.FileType)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 2. 生成存储路径
// 	fileID := generateFileID()
// 	storedName := randomFileName(ext)

// 	// 根据文件类型创建子目录
// 	subDir := req.FileType
// 	if req.FileType == "attachment" {
// 		subDir = "files"
// 	}
// 	saveDir := filepath.Join(s.uploadDir, string(subDir))
// 	if err := os.MkdirAll(saveDir, 0755); err != nil {
// 		return nil, fmt.Errorf("创建目录失败: %w", err)
// 	}

// 	storedPath := filepath.Join(saveDir, storedName)

// 	// 3. 保存文件
// 	srcFile, err := fileHeader.Open()
// 	if err != nil {
// 		return nil, fmt.Errorf("打开上传文件失败: %w", err)
// 	}
// 	defer srcFile.Close()

// 	dstFile, err := os.Create(storedPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("创建目标文件失败: %w", err)
// 	}
// 	defer dstFile.Close()

// 	_, err = io.Copy(dstFile, srcFile)
// 	if err != nil {
// 		return nil, fmt.Errorf("保存文件失败: %w", err)
// 	}

// 	// 4. 如果是图片，处理并获取尺寸
// 	var width, height int
// 	if strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
// 		// 重新打开文件获取尺寸
// 		imgFile, _ := os.Open(storedPath)
// 		defer imgFile.Close()
// 		imgConfig, err := imaging.Decode(imgFile)
// 		if err == nil {
// 			width = imgConfig.Bounds().Dx()
// 			height = imgConfig.Bounds().Dy()
// 		}
// 	}

// 	// 5. 保存到数据库
// 	attachment := &do.Attachment{
// 		FileID:       fileID,
// 		UserID:       userID,
// 		PostID:       req.PostID,
// 		OriginalName: fileHeader.Filename,
// 		StoredName:   storedName,
// 		StoredPath:   storedPath,
// 		Size:         fileHeader.Size,
// 		MimeType:     fileHeader.Header.Get("Content-Type"),
// 		FileType:     req.FileType,
// 		Ext:          ext,
// 		Width:        width,
// 		Height:       height,
// 		Status:       1,
// 	}

// 	if err := s.repo.Create(ctx, attachment); err != nil {
// 		// 上传失败则删除已保存的文件
// 		os.Remove(storedPath)
// 		return nil, fmt.Errorf("保存记录失败: %w", err)
// 	}

// 	// 6. 构建响应
// 	url := fmt.Sprintf("%s%s/%s", s.urlPrefix, subDir, storedName)
// 	return &dto.UploadResponse{
// 		FileID:       fileID,
// 		URL:          url,
// 		OriginalName: fileHeader.Filename,
// 		Size:         fileHeader.Size,
// 		MimeType:     attachment.MimeType,
// 	}, nil
// }

// // 关联帖子
//
//	func (s *service) AssociateWithPost(ctx context.Context, fileID string, postID int64) error {
//		return s.repo.UpdateStatus(ctx, fileID, 1)
//		// 实际应该更新 post_id
//	}
package attachment