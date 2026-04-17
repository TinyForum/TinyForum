package announcement

import (
	"context"
	"time"
	"tiny-forum/internal/model"
	apperrors "tiny-forum/pkg/errors"
)

func (s *announcementService) Publish(ctx context.Context, id uint, userID uint) error {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if announcement.PublishedAt == nil || announcement.PublishedAt.After(time.Now()) {
		return apperrors.ErrInvalidPublishTime
	}
	return s.repo.UpdateStatus(ctx, id, model.AnnouncementStatusPublished)
}

func (s *announcementService) Archive(ctx context.Context, id uint, userID uint) error {
	return s.repo.UpdateStatus(ctx, id, model.AnnouncementStatusArchived)
}

func (s *announcementService) Pin(ctx context.Context, id uint, pinned bool, userID uint) error {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	announcement.IsPinned = pinned
	return s.repo.Update(ctx, announcement)
}
