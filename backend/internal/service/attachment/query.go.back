// package attachment

// import (
// 	"context"
// 	"time"
// 	"tiny-forum/internal/model/dto"
// )

// // func (s *service) ListUserPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[vo.PluginMetaVO], error) {

// // 	return s.pluginSvc.ListPlugins(ctx, queryBO)
// // }

// // 获取文件信息
// func (s *service) GetFile(ctx context.Context, fileID string) (*dto.FileInfo, error) {
// 	att, err := s.repo.GetByFileID(ctx, fileID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &dto.FileInfo{
// 		ID:           att.FileID,
// 		UserID:       att.UserID,
// 		PostID:       att.PostID,
// 		OriginalName: att.OriginalName,
// 		StoredName:   att.StoredName,
// 		StoredPath:   att.StoredPath,
// 		Size:         att.Size,
// 		MimeType:     att.MimeType,
// 		FileType:     att.FileType,
// 		Ext:          att.Ext,
// 		Status:       att.Status,
// 		CreatedAt:    att.CreatedAt.Format(time.RFC3339),
// 	}, nil
// }

// // 获取用户文件列表
// func (s *service) GetUserFiles(ctx context.Context, userID int64, fileType string, page, pageSize int) ([]*dto.FileInfo, int64, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if pageSize < 1 || pageSize > 100 {
// 		pageSize = 20
// 	}
// 	offset := (page - 1) * pageSize

// 	attachments, total, err := s.repo.ListByUser(ctx, userID, fileType, pageSize, offset)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	result := make([]*dto.FileInfo, len(attachments))
// 	for i, att := range attachments {
// 		result[i] = &dto.FileInfo{
// 			ID:           att.FileID,
// 			UserID:       att.UserID,
// 			PostID:       att.PostID,
// 			OriginalName: att.OriginalName,
// 			StoredName:   att.StoredName,
// 			StoredPath:   att.StoredPath,
// 			Size:         att.Size,
// 			MimeType:     att.MimeType,
// 			FileType:     att.FileType,
// 			Ext:          att.Ext,
// 			Status:       att.Status,
// 			CreatedAt:    att.CreatedAt.Format(time.RFC3339),
// 		}
// 	}

// 	return result, total, nil
// }

package attachment