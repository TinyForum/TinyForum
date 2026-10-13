package attachment

import (
	"context"
	"io"
	"mime/multipart"
	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/repository/attachment"
	"tiny-forum/internal/service/upload"
	"tiny-forum/internal/storage"
)

type AttachmentService interface {
	// OpenFile 获取文件内容流，用于直接下载/预览
	OpenFile(ctx context.Context, fileID string) (io.ReadCloser, *dto.FileInfo, error)
	UploadFile(ctx context.Context, userID uint, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest, clientIP string) (*dto.UploadResponse, error) // 上传文件
	GetFile(ctx context.Context, fileID string) (*dto.FileInfo, error)                                                                                               // 获取文件信息
	DeleteFile(ctx context.Context, userID uint, fileID string) error                                                                                                // 删除文件
	GetUserFiles(ctx context.Context, userID uint, fileType string, page, pageSize int) ([]*dto.FileInfo, int64, error)                                              // 获取用户上传的文件列表
	AssociateWithPost(ctx context.Context, fileID string, postID int64) error                                                                                        // 关联附件与帖子
	DeleteByPostID(ctx context.Context, postID int64) error                                                                                                          // 删除帖子的所有附件（物理文件+数据库记录）
}

type service struct {
	repo      attachment.AttachmentRepository
	uploadEng upload.Engine
	storage   storage.StorageDriver
	urlPrefix string
}

func NewAttachmentService(
	repo attachment.AttachmentRepository,
	cfg config.UploadConfig,
	uploadEng upload.Engine,
	store storage.StorageDriver,
) AttachmentService {
	return &service{
		repo:      repo,
		uploadEng: uploadEng,
		storage:   store,
		urlPrefix: cfg.URLPrefix,
	}
}
