// // internal/service/attachment_list_service.go
// package attachment

// import (
// 	"context"
// 	"tiny-forum/internal/model/do"
// 	"tiny-forum/internal/repository/attachment"
// )

// type AttachmentListService interface {
// 	ListUserAttachments(ctx context.Context, userID int64, fileType *do.FileType, page, pageSize int) ([]*do.Attachment, int64, error)
// 	ListPluginAttachments(ctx context.Context, pluginID string, page, pageSize int) ([]*do.Attachment, int64, error)
// }

// type attachmentListServiceImpl struct {
// 	repo attachment.AttachmentRepository
// }

// func NewAttachmentListService(repo attachment.AttachmentRepository) AttachmentListService {
// 	return &attachmentListServiceImpl{repo: repo}
// }

// // func (s *attachmentListServiceImpl) ListUserAttachments(ctx context.Context, userID int64, fileType *do.FileType, page, pageSize int) ([]*do.Attachment, int64, error) {
// // 	status := do.StatusNormal
// // 	query := &repository.AttachmentQuery{
// // 		UserID:   &userID,
// // 		FileType: fileType,
// // 		Status:   &status,
// // 		Page:     page,
// // 		PageSize: pageSize,
// // 	}
// // 	return s.repo.List(ctx, query)
// // }

// func (s *attachmentListServiceImpl) ListPluginAttachments(ctx context.Context, pluginID string, page, pageSize int) ([]*do.Attachment, int64, error) {
// 	status := do.StatusNormal
// 	query := &repository.AttachmentQuery{
// 		PluginID: &pluginID,
// 		Status:   &status,
// 		Page:     page,
// 		PageSize: pageSize,
// 	}
// 	return s.repo.List(ctx, query)
// }

package attachment