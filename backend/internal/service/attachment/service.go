// internal/service/attachment/service.go
package attachment

import (
	"context"
	"errors"
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
	"gorm.io/gorm"
)

type AttachmentService interface {
	UploadFile(ctx context.Context, userID uint, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest, clientIP string) (*dto.UploadResponse, error) // 上传文件
	GetFile(ctx context.Context, fileID string) (*dto.FileInfo, error)                                                                                               // 获取文件信息
	DeleteFile(ctx context.Context, userID uint, fileID string) error                                                                                                // 删除文件
	GetUserFiles(ctx context.Context, userID uint, fileType string, page, pageSize int) ([]*dto.FileInfo, int64, error)                                              // 获取用户上传的文件列表
	AssociateWithPost(ctx context.Context, fileID string, postID int64) error                                                                                        // 关联附件与帖子

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
	// 1. 查询记录（包括已软删除的）
	att, err := s.repo.GetByFileIDUnscoped(ctx, fileID) // 需要新增此方法，不加 deleted_at 过滤
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warnf("文件 %s 不存在", fileID)
			return apperrors.ErrNotFound
		}
		logger.Errorf("查询文件记录失败: %v", err)
		return apperrors.ErrInternalError
	}

	// 2. 如果已经被软删除，幂等返回成功
	if att.DeletedAt.Valid {
		logger.Infof("文件 %s 已被删除，无需重复删除", fileID)
		return nil
	}

	// 3. 权限检查
	if att.UserID != userID {
		logger.Warnf("用户 %d 尝试删除文件 %s，文件所有者是 %d", userID, fileID, att.UserID)
		return apperrors.ErrInsufficientPermission
	}

	// 4. 先删除物理文件，失败则整体失败（避免不一致）
	if err := s.uploadEng.DeleteFile(ctx, att.StoredPath); err != nil {
		logger.Errorf("删除物理文件失败 (path=%s): %v", att.StoredPath, err)
		return apperrors.ErrInternalError // 可自定义：文件删除失败，请稍后重试
	}

	// 5. 软删除数据库记录（如果业务需要软删除）
	if err := s.repo.SoftDelete(ctx, fileID); err != nil {
		logger.Errorf("软删除数据库记录失败: %v", err)
		return apperrors.ErrInternalError
	}

	logger.Infof("用户 %d 成功删除文件 %s", userID, fileID)
	return nil
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
