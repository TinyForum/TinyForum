package attachment

import (
	"context"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type AttachmentRepository interface {
	Create(ctx context.Context, att *do.Attachment) error // 创建文件
	GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error) // 根据文件ID获取文件信息
	Update(ctx context.Context, att *do.Attachment) error // 更新文件信息
	Delete(ctx context.Context, fileID string) error // 删除文件
	ListByUser(ctx context.Context, userID uint, fileType *do.FileType, page, pageSize int) ([]*do.Attachment, int64, error) // 根据用户ID获取文件列表
	FindDuplicate(ctx context.Context, fileHash string, fileType do.FileType) (*do.Attachment, error) // 查找重复文件
	AssociateWithPost(ctx context.Context, fileID string, postID int64) error // 将文件与帖子关联
	CheckFileExist(ctx context.Context, fileID string) bool // 检查文件是否存在
	SoftDelete(ctx context.Context, fileID string) error // 软删除
	GetByFileIDUnscoped(ctx context.Context, fileID string) (*do.Attachment, error) // 获取文件信息（包括软删除的）
}

type attachmentRepo struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepo{db: db}
}

func (r *attachmentRepo) Create(ctx context.Context, att *do.Attachment) error {
	return r.db.WithContext(ctx).Create(att).Error
}

// GetByFileIDUnscoped 即使软删除也能查到
func (r *attachmentRepo) GetByFileIDUnscoped(ctx context.Context, fileID string) (*do.Attachment, error) {
  var att do.Attachment
    err := r.db.Unscoped().Where("file_id = ?", fileID).First(&att).Error
    return &att, err
}

// SoftDelete 软删除（如果 GORM 模型已定义 DeletedAt）
func (r *attachmentRepo) SoftDelete(ctx context.Context, fileID string) error {
	var att do.Attachment
    return r.db.Where("file_id = ?", fileID).Delete(&att).Error
}

// 检查文件是否存在
func (r *attachmentRepo) CheckFileExist(ctx context.Context, fileID string) bool {
    return r.db.WithContext(ctx).Where("file_id = ? AND deleted_at IS NULL", fileID).First(&do.Attachment{}).Error == nil
}
// 获取文件信息
func (r *attachmentRepo) GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error) {
	var att do.Attachment
	err := r.db.WithContext(ctx).Where("file_id = ?", fileID).First(&att).Error
	return &att, err
}

func (r *attachmentRepo) Update(ctx context.Context, att *do.Attachment) error {
	return r.db.WithContext(ctx).Save(att).Error
}

func (r *attachmentRepo) Delete(ctx context.Context, fileID string) error {
	return r.db.WithContext(ctx).Where("file_id = ?", fileID).Delete(&do.Attachment{}).Error
}

func (r *attachmentRepo) ListByUser(ctx context.Context, userID uint, fileType *do.FileType, page, pageSize int) ([]*do.Attachment, int64, error) {
	query := r.db.WithContext(ctx).Model(&do.Attachment{}).Where("user_id = ? AND status = ?", userID, do.StatusNormal)
	if fileType != nil {
		query = query.Where("file_type = ?", *fileType)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*do.Attachment
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&list).Error
	return list, total, err
}

func (r *attachmentRepo) FindDuplicate(ctx context.Context, fileHash string, fileType do.FileType) (*do.Attachment, error) {
	var att do.Attachment
	err := r.db.WithContext(ctx).
		Where("file_hash = ? AND file_type = ? AND status = ?", fileHash, fileType, do.StatusNormal).
		First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *attachmentRepo) AssociateWithPost(ctx context.Context, fileID string, postID int64) error {
	return r.db.WithContext(ctx).Model(&do.Attachment{}).
		Where("file_id = ?", fileID).
		Update("post_id", postID).Error
}