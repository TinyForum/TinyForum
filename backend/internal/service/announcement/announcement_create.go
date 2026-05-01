package announcement

import (
	"context"
	"time"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/po"
)

func (s *announcementService) Create(ctx context.Context, req *dto.CreateAnnouncementRequest, userID uint) (*po.Announcement, error) {
	if err := s.validateTime(req.PublishedAt, req.ExpiredAt); err != nil {
		return nil, err
	}
	now := time.Now()
	announcement := &po.Announcement{
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
		Status:      po.AnnouncementStatusDraft,
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
