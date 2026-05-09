package attachment

import (
	"context"
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// AttachmentRepository 定义了附件存储的数据库操作接口
type AttachmentRepository interface {
	Create(ctx context.Context, att *do.Attachment) error // 创建新附件记录
	GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error)          // 仅查询未软删除的记录
	GetByFileIDUnscoped(ctx context.Context, fileID string) (*do.Attachment, error) // 包含已软删除的记录
	Update(ctx context.Context, att *do.Attachment) error
	Delete(ctx context.Context, fileID string) error          // 硬删除（物理删除）
	SoftDelete(ctx context.Context, fileID string) error      // 软删除（设置 deleted_at）
	CheckFileExist(ctx context.Context, fileID string) bool   // 检查未软删除的记录是否存在
	ListByUser(ctx context.Context, userID uint, fileType *do.FileType, page, pageSize int) ([]*do.Attachment, int64, error) // 根据用户ID获取附件列表
	FindDuplicate(ctx context.Context, fileHash string, fileType do.FileType) (*do.Attachment, error) // 查找重复文件
	AssociateWithPost(ctx context.Context, fileID string, postID int64) error // 将附件关联到帖子
}

type attachmentRepo struct {
	db *gorm.DB
}

// NewAttachmentRepository 创建 AttachmentRepository 实例
func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepo{db: db}
}

// Create 创建新附件记录
func (r *attachmentRepo) Create(ctx context.Context, att *do.Attachment) error {
	return r.db.WithContext(ctx).Create(att).Error
}

// GetByFileID 根据 fileID 获取未软删除的附件信息
func (r *attachmentRepo) GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error) {
	var att do.Attachment
	err := r.db.WithContext(ctx).
		Where("file_id = ? AND deleted_at IS NULL", fileID).
		First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

// GetByFileIDUnscoped 根据 fileID 获取附件信息（包含已软删除的记录）
func (r *attachmentRepo) GetByFileIDUnscoped(ctx context.Context, fileID string) (*do.Attachment, error) {
	var att do.Attachment
	err := r.db.WithContext(ctx).Unscoped().
		Where("file_id = ?", fileID).
		First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

// Update 更新附件记录（仅更新非零值字段）
func (r *attachmentRepo) Update(ctx context.Context, att *do.Attachment) error {
	// 使用 Updates 而非 Save，避免零值覆盖问题
	return r.db.WithContext(ctx).Model(&do.Attachment{}).
		Where("file_id = ?", att.FileID).
		Updates(att).Error
}

// Delete 硬删除附件（物理删除，同时从数据库移除记录）
func (r *attachmentRepo) Delete(ctx context.Context, fileID string) error {
	return r.db.WithContext(ctx).Unscoped().
		Where("file_id = ?", fileID).
		Delete(&do.Attachment{}).Error
}

// SoftDelete 软删除附件（设置 deleted_at 字段）
func (r *attachmentRepo) SoftDelete(ctx context.Context, fileID string) error {
	return r.db.WithContext(ctx).
		Where("file_id = ?", fileID).
		Delete(&do.Attachment{}).Error // GORM 会自动设置 deleted_at 如果模型包含 gorm.DeletedAt
}

// CheckFileExist 检查未软删除的附件是否存在
func (r *attachmentRepo) CheckFileExist(ctx context.Context, fileID string) bool {
	var count int64
	r.db.WithContext(ctx).Model(&do.Attachment{}).
		Where("file_id = ? AND deleted_at IS NULL", fileID).
		Count(&count)
	return count > 0
}

// ListByUser 分页获取用户的附件列表（仅未软删除）
func (r *attachmentRepo) ListByUser(ctx context.Context, userID uint, fileType *do.FileType, page, pageSize int) ([]*do.Attachment, int64, error) {
	query := r.db.WithContext(ctx).Model(&do.Attachment{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	if fileType != nil && *fileType != "" {
		query = query.Where("file_type = ?", *fileType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*do.Attachment
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&list).Error
	return list, total, err
}

// FindDuplicate 根据文件 hash 和类型查找重复的未删除附件
func (r *attachmentRepo) FindDuplicate(ctx context.Context, fileHash string, fileType do.FileType) (*do.Attachment, error) {
	var att do.Attachment
	err := r.db.WithContext(ctx).
		Where("file_hash = ? AND file_type = ? AND deleted_at IS NULL", fileHash, fileType).
		First(&att).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 未找到重复记录不是错误，返回 nil
		}
		return nil, err
	}
	return &att, nil
}

// AssociateWithPost 将附件关联到指定帖子（更新 post_id）
func (r *attachmentRepo) AssociateWithPost(ctx context.Context, fileID string, postID int64) error {
	return r.db.WithContext(ctx).Model(&do.Attachment{}).
		Where("file_id = ? AND deleted_at IS NULL", fileID).
		Update("post_id", postID).Error
}