package announcement

import (
	"context"
	"time"
	"tiny-forum/internal/model"
)

func (s *announcementService) Create(ctx context.Context, req *CreateAnnouncementRequest, userID uint) (*model.Announcement, error) {
	if err := s.validateTime(req.PublishedAt, req.ExpiredAt); err != nil {
		return nil, err
	}
	now := time.Now()
	announcement := &model.Announcement{
		Title:       req.Title,
		Content:     req.Content,
		Summary:     req.Summary,
		Cover:       req.Cover,
		Type:        req.Type,
		IsPinned:    req.IsPinned,
		IsGlobal:    req.IsGlobal,
		BoardID:     req.BoardID,
		PublishedAt: req.PublishedAt,
		ExpiredAt:   req.ExpiredAt,
		Status:      model.AnnouncementStatusDraft,
		ViewCount:   0,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}
	if announcement.PublishedAt == nil {
		announcement.PublishedAt = &now
	}
	if err := s.repo.Create(ctx, announcement); err != nil {
		return nil, err
	}
	return announcement, nil
}
