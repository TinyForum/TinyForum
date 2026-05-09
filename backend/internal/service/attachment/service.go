// internal/service/attachment/service.go
package attachment

import (
	"context"
	"fmt"
	"mime/multipart"

	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	attachmentRepo "tiny-forum/internal/repository/attachment"
	"tiny-forum/internal/service/upload"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

	"github.com/google/uuid"
)

type AttachmentService interface {
	UploadFile(ctx context.Context, userID uint, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest, clientIP string) (*dto.UploadResponse, error)
	GetFile(ctx context.Context, fileID string) (*dto.FileInfo, error)
	DeleteFile(ctx context.Context, userID uint, fileID string) error
	GetUserFiles(ctx context.Context, userID uint, fileType string, page, pageSize int) ([]*dto.FileInfo, int64, error)
	AssociateWithPost(ctx context.Context, fileID string, postID int64) error
}

type service struct {
	repo      attachmentRepo.AttachmentRepository
	uploadEng upload.Engine
	urlPrefix string
}

func NewAttachmentService(
	repo attachmentRepo.AttachmentRepository,
	cfg config.UploadConfig,
	uploadEng upload.Engine,
) AttachmentService {
	return &service{
		repo:      repo,
		uploadEng: uploadEng,
		urlPrefix: cfg.URLPrefix,
	}
}

func (s *service) UploadFile(ctx context.Context, userID uint, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest, clientIP string) (*dto.UploadResponse, error) {
	// 1. 先调用上传引擎存储文件（不写入数据库）
	uploadReq := &upload.UploadRequest{
		UserID:   userID,
		File:     fileHeader,
		FileType: do.FileType(req.Type),
		PostID:   req.PostID,
		ReplyID:  req.ReplyID,
		ClientIP: clientIP,
	}
	result, err := s.uploadEng.Upload(ctx, uploadReq)
	if err != nil {
		return nil, err
	}

	// 2. 检查是否已存在相同哈希的文件（去重）
	dup, _ := s.repo.FindDuplicate(ctx, result.FileHash, do.FileType(req.Type))
	if dup != nil {
		// 已存在，返回已有文件的信息（物理文件已存在，无需重复保存）
		return &dto.UploadResponse{
			FileID: dup.FileID,
			URL:    s.urlPrefix + "/" + dup.StoredPath,
		}, nil
	}

	// 3. 保存元数据到数据库
	meta := &do.Attachment{
		FileID:       uuid.New().String(),
		UserID:       userID,
		PluginID:     "", // 普通上传没有 plugin_id
		PostID:       req.PostID,
		ReplyID:      req.ReplyID,
		OriginalName: result.OriginalName,
		StoredName:   result.StoredName,
		StoredPath:   result.StoredPath,
		Size:         result.Size,
		FileType:     do.FileType(req.Type),
		MimeType:     result.MimeType,
		MimeMajor:    result.MimeMajor,
		Ext:          result.Ext,
		Status:       do.StatusNormal, // 直接标记为正常（或者你希望先临时再异步更新，也可以）
		UploadIP:     clientIP,
		FileHash:     result.FileHash,
	}
	if err := s.repo.Create(ctx, meta); err != nil {
		// 如果数据库写入失败，尝试删除已存储的物理文件
				_ = s.uploadEng.DeleteFile(ctx, result.StoredPath)
		return nil, fmt.Errorf("save record: %w", err)
	}

	return &dto.UploadResponse{
		FileID: meta.FileID,
		URL:    s.urlPrefix + "/" + meta.StoredPath,
	}, nil
}

func (s *service) GetFile(ctx context.Context, fileID string) (*dto.FileInfo, error) {
	att, err := s.repo.GetByFileID(ctx, fileID)
	if err != nil {
		return nil, err
	}
	return &dto.FileInfo{
		FileID:       att.FileID,
		OriginalName: att.OriginalName,
		Size:         att.Size,
		FileType:     string(att.FileType),
		MimeType:     att.MimeType,
		URL:          s.urlPrefix + "/" + att.StoredPath,
		// CreatedAt:    att.CreatedAt,
	}, nil
}

func (s *service) DeleteFile(ctx context.Context, userID uint, fileID string) error {
	att, err := s.repo.GetByFileID(ctx, fileID)
	if err != nil {
		logger.Errorf("查询失败: ", err)
		return err
	}
	if att.UserID != userID {
		logger.Errorf("用户 %d 尝试删除文件 %s，该文件所有者是 %d", userID, fileID, att.UserID)
		return apperrors.ErrInsufficientPermission
	}
	// 删除物理文件
	if err := s.uploadEng.DeleteFile(ctx, att.StoredPath); err != nil {
		logger.Errorf("删除物理文件失败: %v", err)
	}
	if err := s.repo.Delete(ctx, fileID); err != nil {
	    logger.Errorf("删除数据库文件失败: ", err)
	}

	// 删除数据库记录
	return err
}
func (s *service) GetUserFiles(ctx context.Context, userID uint, fileTypeStr string, page, pageSize int) ([]*dto.FileInfo, int64, error) {
	var fileType *do.FileType
	if fileTypeStr != "" {
		ft := do.FileType(fileTypeStr)
		fileType = &ft
	}
	attachments, total, err := s.repo.ListByUser(ctx, userID, fileType, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*dto.FileInfo, len(attachments))
	for i, att := range attachments {
		result[i] = &dto.FileInfo{
			FileID:       att.FileID,
			OriginalName: att.OriginalName,
			Size:         att.Size,
			FileType:     string(att.FileType),
			MimeType:     att.MimeType,
			URL:          s.urlPrefix + "/" + att.StoredPath,
			CreatedAt:    att.CreatedAt,
		}
	}
	return result, total, nil
}

func (s *service) AssociateWithPost(ctx context.Context, fileID string, postID int64) error {
	return s.repo.AssociateWithPost(ctx, fileID, postID)
}