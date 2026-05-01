package upload

import (
	"context"
	"mime/multipart"
	"os"
	"strings"
	"tiny-forum/config"
	"tiny-forum/internal/model/dto"
	uploadRepo "tiny-forum/internal/repository/upload"
)

type UploadService interface {
	UploadFile(ctx context.Context, userID int64, fileHeader *multipart.FileHeader, req *dto.UploadRequest) (*dto.UploadResponse, error)
	GetFile(ctx context.Context, fileID string) (*dto.FileInfo, error)
	DeleteFile(ctx context.Context, userID int64, fileID string) error
	GetUserFiles(ctx context.Context, userID int64, fileType string, page, pageSize int) ([]*dto.FileInfo, int64, error)
	AssociateWithPost(ctx context.Context, fileID string, postID int64) error
}

type service struct {
	repo       uploadRepo.UploadRepository
	uploadDir  string
	urlPrefix  string
	maxSize    int64
	allowedExt map[string]bool
}

func NewUploadService(repo uploadRepo.UploadRepository, cfg config.UploadConfig) UploadService {
	allowedMap := make(map[string]bool)
	for _, ext := range cfg.AllowedExt {
		allowedMap[strings.ToLower(ext)] = true
	}

	// 确保上传目录存在
	os.MkdirAll(cfg.UploadDir, 0755)

	return &service{
		repo:       repo,
		uploadDir:  cfg.UploadDir,
		urlPrefix:  cfg.URLPrefix,
		maxSize:    cfg.MaxSize,
		allowedExt: allowedMap,
	}
}
