package attachment

import (
	"context"
	"mime/multipart"
	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/repository/attachment"
	"tiny-forum/internal/service/upload"
)

type AttachmentService interface {
	UploadFile(ctx context.Context, userID uint, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest, clientIP string) (*dto.UploadResponse, error) // 上传文件
	GetFile(ctx context.Context, fileID string) (*dto.FileInfo, error)                                                                                               // 获取文件信息
	DeleteFile(ctx context.Context, userID uint, fileID string) error                                                                                                // 删除文件
	GetUserFiles(ctx context.Context, userID uint, fileType string, page, pageSize int) ([]*dto.FileInfo, int64, error)                                              // 获取用户上传的文件列表
	AssociateWithPost(ctx context.Context, fileID string, postID int64) error                                                                                        // 关联附件与帖子

}

type service struct {
	repo      attachment.AttachmentRepository
	uploadEng upload.Engine
	urlPrefix string
}

func NewAttachmentService(
	repo attachment.AttachmentRepository,
	cfg config.UploadConfig,
	uploadEng upload.Engine,
) AttachmentService {
	return &service{
		repo:      repo,
		uploadEng: uploadEng,
		urlPrefix: cfg.URLPrefix,
	}
}
