package announcement

import (
	"context"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *announcementRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&do.Announcement{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *announcementRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&do.Announcement{}, ids).Error
}

func (r *announcementRepository) UpdateStatus(ctx context.Context, id uint, status do.AnnouncementStatus) error {
	return r.db.WithContext(ctx).
		Model(&do.Announcement{}).
		Where("id = ?", id).
		Update("status", status).Error
}
