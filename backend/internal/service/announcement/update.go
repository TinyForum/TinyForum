package announcement

import (
	"context"
	"errors"
	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

func (s *announcementService) Update(ctx context.Context, id uint, req *request.UpdateAnnouncement, userID uint) error {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrAnnouncementNotFound
		}
		return err
	}
	if req.Title != nil {
		announcement.Title = *req.Title
	}
	if req.Content != nil {
		announcement.Content = *req.Content
	}
	if req.Summary != nil {
		announcement.Summary = *req.Summary
	}
	if req.Cover != nil {
		announcement.Cover = *req.Cover
	}
	if req.Type != nil {
		announcement.Type = req.Type
	}
	if req.IsPinned != nil {
		announcement.IsPinned = *req.IsPinned
	}
	if req.IsGlobal != nil {
		announcement.IsGlobal = *req.IsGlobal
	}
	if req.BoardID != nil {
		announcement.BoardID = req.BoardID
	}
	if req.PublishedAt != nil {
		announcement.PublishedAt = req.PublishedAt
	}
	if req.ExpiredAt != nil {
		announcement.ExpiredAt = req.ExpiredAt
	}
	if err := s.validateTime(announcement.PublishedAt, announcement.ExpiredAt); err != nil {
		return err
	}
	announcement.UpdatedBy = userID
	return s.repo.Update(ctx, announcement)
}
