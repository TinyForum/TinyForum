package attachment

import (
	"context"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// AttachmentRepository 定义了附件存储的数据库操作接口
type AttachmentRepository interface {
	Create(ctx context.Context, att *do.Attachment) error                                                                    // 创建新附件记录
	GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error)                                                  // 仅查询未软删除的记录
	GetByFileIDUnscoped(ctx context.Context, fileID string) (*do.Attachment, error)                                          // 包含已软删除的记录
	Update(ctx context.Context, att *do.Attachment) error                                                                    // 更新附件记录（仅更新非零值字段）
	Delete(ctx context.Context, fileID string) error                                                                         // 硬删除（物理删除）
	SoftDelete(ctx context.Context, fileID string) error                                                                     // 软删除（设置 deleted_at）
	CheckFileExist(ctx context.Context, fileID string) bool                                                                  // 检查未软删除的记录是否存在
	ListByUser(ctx context.Context, userID uint, fileType *do.FileType, page, pageSize int) ([]*do.Attachment, int64, error) // 根据用户ID获取附件列表
	FindDuplicate(ctx context.Context, fileHash string, fileType do.FileType) (*do.Attachment, error)                        // 查找重复文件
	AssociateWithPost(ctx context.Context, fileID string, postID int64) error                                                // 将附件关联到帖子
}

type attachmentRepo struct {
	db *gorm.DB
}

// NewAttachmentRepository 创建 AttachmentRepository 实例
func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepo{db: db}
}
