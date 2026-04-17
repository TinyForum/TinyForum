package announcement

import (
	"context"
	"tiny-forum/internal/model"
)

func (r *announcementRepository) Create(ctx context.Context, announcement *model.Announcement) error {
	return r.db.WithContext(ctx).Create(announcement).Error
}

func (r *announcementRepository) Update(ctx context.Context, announcement *model.Announcement) error {
	return r.db.WithContext(ctx).Save(announcement).Error
}

func (r *announcementRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Announcement{}, id).Error
}

func (r *announcementRepository) GetByID(ctx context.Context, id uint) (*model.Announcement, error) {
	var announcement model.Announcement
	err := r.db.WithContext(ctx).
		Preload("Board").
		Preload("Creator").
		First(&announcement, id).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}
