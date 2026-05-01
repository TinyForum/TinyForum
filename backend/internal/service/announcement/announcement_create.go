package announcement

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

func (s *announcementService) Create(ctx context.Context, req *request.CreateAnnouncement, userID uint) (*do.Announcement, error) {
	if err := s.validateTime(req.PublishedAt, req.ExpiredAt); err != nil {
		return nil, err
	}
	now := time.Now()
	announcement := &do.Announcement{
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
		Status:      do.AnnouncementStatusDraft,
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
