package upload

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
)

func (r *repository) Create(ctx context.Context, attachment *do.Attachment) error {
	return r.db.WithContext(ctx).Create(attachment).Error
}

func (r *repository) GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error) {
	var att do.Attachment
	err := r.db.WithContext(ctx).Where("file_id = ?", fileID).First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *repository) GetByPostID(ctx context.Context, postID int64, limit, offset int) ([]*do.Attachment, error) {
	var list []*do.Attachment
	err := r.db.WithContext(ctx).
		Where("post_id = ? AND status = 1", postID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&list).Error
	return list, err
}

func (r *repository) Delete(ctx context.Context, fileID string) error {
	return r.db.WithContext(ctx).
		Model(&do.Attachment{}).
		Where("file_id = ?", fileID).
		Update("status", 2).Error
}

func (r *repository) UpdateStatus(ctx context.Context, fileID string, status int) error {
	return r.db.WithContext(ctx).
		Model(&do.Attachment{}).
		Where("file_id = ?", fileID).
		Update("status", status).Error
}

func (r *repository) ListByUser(ctx context.Context, userID int64, fileType string, limit, offset int) ([]*do.Attachment, int64, error) {
	var list []*do.Attachment
	var total int64

	query := r.db.WithContext(ctx).Model(&do.Attachment{}).Where("user_id = ? AND status = 1", userID)
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *repository) DeleteUnusedTemp(ctx context.Context, beforeTime time.Time) (int64, error) {
	result := r.db.WithContext(ctx).
		Where("status = 0 AND created_at < ?", beforeTime).
		Delete(&do.Attachment{})
	return result.RowsAffected, result.Error
}
