// package attachment

// import (
// 	"context"
// 	"time"
// 	"tiny-forum/internal/model/do"

// 	"gorm.io/gorm"
// )

// type AttachmentRepository interface {
// 	Create(ctx context.Context, attachment *do.Attachment) error // 创建
// 	GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error)
// 	GetByPostID(ctx context.Context, postID int64, limit, offset int) ([]*do.Attachment, error)
// 	Delete(ctx context.Context, fileID string) error
// 	Update(ctx context.Context, att *do.Attachment) error
// 	UpdateStatus(ctx context.Context, fileID string, status int) error
// 	ListByUser(ctx context.Context, userID int64, fileType string, limit, offset int) ([]*do.Attachment, int64, error)
// 	DeleteUnusedTemp(ctx context.Context, beforeTime time.Time) (int64, error)
// 	 FindDuplicate(ctx context.Context, fileHash string, fileType do.FileType) (*do.Attachment, error) // 新增
// }

// type repository struct {
// 	db *gorm.DB
// }

// // NewAttachmentRepository 创建一个新的附件仓库实例
// // 参数:
// //   db: GORM数据库连接对象，用于数据库操作
// // 返回值:
// //   AttachmentRepository: 附件仓库接口的实现，这里返回的是repository结构体的指针
//
//	func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
//		return &repository{db: db}
//	}
package attachment