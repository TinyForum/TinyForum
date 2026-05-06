package upload

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/logger"
)

// 生成唯一文件ID
func generateFileID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// 生成随机文件名
func randomFileName(ext string) string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), hex.EncodeToString(b), ext)
}

// 检查文件类型
func (s *service) validateFile(header *multipart.FileHeader, fileType do.FileType) (string, error) {
	// 检查大小
	if header.Size > s.maxSize {
		return "", fmt.Errorf("文件过大，最大允许 %d MB", s.maxSize/(1024*1024))
	}

	// 检查扩展名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !s.allowedExt[ext] {
		return "", fmt.Errorf("不支持的文件类型: %s", ext)
	}

	// 嗅探真实MIME类型
	file, err := header.Open()
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	sniffBuf := make([]byte, 512)
	_, err = file.Read(sniffBuf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("读取文件头失败: %w", err)
	}

	mime := http.DetectContentType(sniffBuf)

	// 根据类型额外验证
	switch fileType {
	case "avatar":
		if !strings.HasPrefix(mime, "image/") {
			return "", fmt.Errorf("头像必须是图片格式")
		}
	case "post_image":
		if !strings.HasPrefix(mime, "image/") {
			return "", fmt.Errorf("帖子图片必须是图片格式")
		}
	case "comment_attachment":
		// 附件允许更多类型
		allowedMimes := []string{"image/", "application/pdf", "application/zip"}
		allowed := false
		for _, allowedMime := range allowedMimes {
			if strings.HasPrefix(mime, allowedMime) {
				allowed = true
				break
			}
		}
		if !allowed {
			return "", fmt.Errorf("不支持的附件类型")
		}
	}

	return ext, nil
}

// 处理图片（缩放、格式转换）
func processImage(srcPath, dstPath string, maxWidth, maxHeight int) error {
	src, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}

	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 如果图片尺寸超过限制，进行缩放
	if width > maxWidth || height > maxHeight {
		src = imaging.Fit(src, maxWidth, maxHeight, imaging.Lanczos)
	}

	return imaging.Save(src, dstPath)
}

// 主上传方法
func (s *service) UploadFile(ctx context.Context, userID int64, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest) (*dto.UploadResponse, error) {
	// 1. 验证文件
	ext, err := s.validateFile(fileHeader, req.FileType)
	if err != nil {
		return nil, err
	}

	// 2. 生成存储路径
	fileID := generateFileID()
	storedName := randomFileName(ext)

	// 根据文件类型创建子目录
	subDir := req.FileType
	if req.FileType == "attachment" {
		subDir = "files"
	}
	saveDir := filepath.Join(s.uploadDir, string(subDir))
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	storedPath := filepath.Join(saveDir, storedName)

	// 3. 保存文件
	srcFile, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("打开上传文件失败: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(storedPath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 4. 如果是图片，处理并获取尺寸
	var width, height int
	if strings.HasPrefix(fileHeader.Header.Get("Content-Type"), "image/") {
		// 重新打开文件获取尺寸
		imgFile, _ := os.Open(storedPath)
		defer imgFile.Close()
		imgConfig, err := imaging.Decode(imgFile)
		if err == nil {
			width = imgConfig.Bounds().Dx()
			height = imgConfig.Bounds().Dy()
		}
	}

	// 5. 保存到数据库
	attachment := &do.Attachment{
		FileID:       fileID,
		UserID:       userID,
		PostID:       req.PostID,
		OriginalName: fileHeader.Filename,
		StoredName:   storedName,
		StoredPath:   storedPath,
		Size:         fileHeader.Size,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		FileType:     req.FileType,
		Ext:          ext,
		Width:        width,
		Height:       height,
		Status:       1,
	}

	if err := s.repo.Create(ctx, attachment); err != nil {
		// 上传失败则删除已保存的文件
		os.Remove(storedPath)
		return nil, fmt.Errorf("保存记录失败: %w", err)
	}

	// 6. 构建响应
	url := fmt.Sprintf("%s%s/%s", s.urlPrefix, subDir, storedName)
	return &dto.UploadResponse{
		FileID:       fileID,
		URL:          url,
		OriginalName: fileHeader.Filename,
		Size:         fileHeader.Size,
		MimeType:     attachment.MimeType,
	}, nil
}

// 获取文件信息
func (s *service) GetFile(ctx context.Context, fileID string) (*dto.FileInfo, error) {
	att, err := s.repo.GetByFileID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	return &dto.FileInfo{
		ID:           att.FileID,
		UserID:       att.UserID,
		PostID:       att.PostID,
		OriginalName: att.OriginalName,
		StoredName:   att.StoredName,
		StoredPath:   att.StoredPath,
		Size:         att.Size,
		MimeType:     att.MimeType,
		FileType:     att.FileType,
		Ext:          att.Ext,
		Status:       att.Status,
		CreatedAt:    att.CreatedAt.Format(time.RFC3339),
	}, nil
}

// 删除文件
func (s *service) DeleteFile(ctx context.Context, userID int64, fileID string) error {
	att, err := s.repo.GetByFileID(ctx, fileID)
	if err != nil {
		return err
	}

	// 检查权限：只有文件所有者或管理员可以删除
	if att.UserID != userID {
		// TODO: 添加管理员权限检查
		return fmt.Errorf("无权删除此文件")
	}

	// 删除数据库记录
	if err := s.repo.Delete(ctx, fileID); err != nil {
		return err
	}

	// 删除物理文件（可选：保留一段时间或异步删除）
	if err := os.Remove(att.StoredPath); err != nil {
		logger.Warnf("删除物理文件失败: %v", err)
	}

	return nil
}

// 获取用户文件列表
func (s *service) GetUserFiles(ctx context.Context, userID int64, fileType string, page, pageSize int) ([]*dto.FileInfo, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	attachments, total, err := s.repo.ListByUser(ctx, userID, fileType, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*dto.FileInfo, len(attachments))
	for i, att := range attachments {
		result[i] = &dto.FileInfo{
			ID:           att.FileID,
			UserID:       att.UserID,
			PostID:       att.PostID,
			OriginalName: att.OriginalName,
			StoredName:   att.StoredName,
			StoredPath:   att.StoredPath,
			Size:         att.Size,
			MimeType:     att.MimeType,
			FileType:     att.FileType,
			Ext:          att.Ext,
			Status:       att.Status,
			CreatedAt:    att.CreatedAt.Format(time.RFC3339),
		}
	}

	return result, total, nil
}

// 关联帖子
func (s *service) AssociateWithPost(ctx context.Context, fileID string, postID int64) error {
	return s.repo.UpdateStatus(ctx, fileID, 1)
	// 实际应该更新 post_id
}
