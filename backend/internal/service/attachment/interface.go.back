// package attachment

// import (
// 	"context"
// 	"os"
// 	"strings"
// 	"tiny-forum/internal/infra/config"
// 	"tiny-forum/internal/model/dto"
// 	attachmentRepo "tiny-forum/internal/repository/attachment"
// )

// type AttachmentService interface {
// 	// UploadFile(ctx context.Context, userID int64, fileHeader *multipart.FileHeader, req *request.UploadPostFileRequest) (*dto.UploadResponse, error)
// 	// UploadPlugin(ctx context.Context, request bo.PluginUpdateBO) (*dto.UploadResponse, error)
// 	GetFile(ctx context.Context, fileID string) (*dto.FileInfo, error)
// 	DeleteFile(ctx context.Context, userID int64, fileID string) error
// 	GetUserFiles(ctx context.Context, userID int64, fileType string, page, pageSize int) ([]*dto.FileInfo, int64, error)
// 	AssociateWithPost(ctx context.Context, fileID string, postID int64) error
// 	// ListUserPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[vo.PluginMetaVO], error)

// }

// type service struct {
// 	repo       attachmentRepo.AttachmentRepository
// 	uploadDir  string
// 	urlPrefix  string
// 	maxSize    int64
// 	allowedExt map[string]bool
// 	// pluginSvc  plugin.PluginService
// }

// func NewAttachmentService(
// 	repo attachmentRepo.AttachmentRepository,
// 	cfg config.UploadConfig,
// 	// plugin plugin.PluginService,
// ) AttachmentService {
// 	allowedMap := make(map[string]bool)
// 	for _, ext := range cfg.AllowedExt {
// 		allowedMap[strings.ToLower(ext)] = true
// 	}

// 	// 确保上传目录存在
// 	os.MkdirAll(cfg.UploadDir, 0755)

//		return &service{
//			repo:       repo,
//			uploadDir:  cfg.UploadDir,
//			urlPrefix:  cfg.URLPrefix,
//			maxSize:    cfg.MaxSize,
//			allowedExt: allowedMap,
//			// pluginSvc:  plugin,
//		}
//	}
package attachment