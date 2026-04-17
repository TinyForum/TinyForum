package announcement

import (
	"context"
	"time"
	"tiny-forum/internal/model"
)

func (r *announcementRepository) List(ctx context.Context, req *AnnouncementListRequest) ([]model.Announcement, int64, error) {
	var announcements []model.Announcement
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Announcement{})

	if req.BoardID != nil {
		query = query.Where("board_id = ? OR (board_id IS NULL AND is_global = ?)", *req.BoardID, true)
	}
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.Status != nil {
		if *req.Status != StatusAll {
			query = query.Where("status = ?", *req.Status)
		}
	}
	if req.IsPinned != nil {
		query = query.Where("is_pinned = ?", *req.IsPinned)
	}
	if req.IsGlobal != nil {
		query = query.Where("is_global = ?", *req.IsGlobal)
	}
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.StartTime != nil {
		query = query.Where("published_at >= ?", req.StartTime)
	}
	if req.EndTime != nil {
		query = query.Where("published_at <= ?", req.EndTime)
	}

	shouldFilterTime := req.Status == nil ||
		(req.Status != nil && *req.Status == StatusPublished)

	if shouldFilterTime {
		query = query.Where("published_at <= ?", time.Now())
		query = query.Where("expired_at IS NULL OR expired_at > ?", time.Now())
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	err := query.
		Order("is_pinned DESC").
		Order("published_at DESC").
		Preload("Board").
		Preload("Creator").
		Offset(offset).
		Limit(req.PageSize).
		Find(&announcements).Error

	return announcements, total, err
}

func (r *announcementRepository) GetPinned(ctx context.Context, boardID *uint) ([]model.Announcement, error) {
	var announcements []model.Announcement
	query := r.db.WithContext(ctx).
		Where("is_pinned = ?", true).
		Where("status = ?", model.AnnouncementStatusPublished).
		Where("published_at <= ?", time.Now()).
		Where("expired_at IS NULL OR expired_at > ?", time.Now())

	if boardID != nil {
		query = query.Where("board_id = ? OR (board_id IS NULL AND is_global = ?)", *boardID, true)
	} else {
		query = query.Where("is_global = ?", true)
	}

	err := query.Order("published_at DESC").
		Preload("Board").
		Find(&announcements).Error

	return announcements, err
}
