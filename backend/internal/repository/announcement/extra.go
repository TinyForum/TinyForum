package announcement

import (
	"context"
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

func (r *announcementRepository) IncrementViewCount(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&po.Announcement{}).
		Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *announcementRepository) BatchDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&po.Announcement{}, ids).Error
}

func (r *announcementRepository) UpdateStatus(ctx context.Context, id uint, status po.AnnouncementStatus) error {
	return r.db.WithContext(ctx).
		Model(&po.Announcement{}).
		Where("id = ?", id).
		Update("status", status).Error
}
