package attachment

import (
	"context"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type AttachmentRepository interface {
	Create(ctx context.Context, att *do.Attachment) error
	GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error)
	Update(ctx context.Context, att *do.Attachment) error
	Delete(ctx context.Context, fileID string) error
	ListByUser(ctx context.Context, userID int64, fileType *do.FileType, page, pageSize int) ([]*do.Attachment, int64, error)
	FindDuplicate(ctx context.Context, fileHash string, fileType do.FileType) (*do.Attachment, error)
	AssociateWithPost(ctx context.Context, fileID string, postID int64) error
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

func (r *attachmentRepo) ListByUser(ctx context.Context, userID int64, fileType *do.FileType, page, pageSize int) ([]*do.Attachment, int64, error) {
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