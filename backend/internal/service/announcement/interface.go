package announcement

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	announcementRepo "tiny-forum/internal/repository/announcement"
	apperrors "tiny-forum/pkg/errors"
)

type AnnouncementService interface {
	Create(ctx context.Context, req *request.CreateAnnouncement, userID uint) (*do.Announcement, error)
	Update(ctx context.Context, id uint, req *request.UpdateAnnouncement, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
	GetByID(ctx context.Context, id uint) (*do.Announcement, error)
	List(ctx context.Context, req *request.ListAnnouncements) (*vo.ListAnnouncements, error)
	GetPinned(ctx context.Context, boardID *uint) ([]do.Announcement, error)
	Publish(ctx context.Context, id uint, userID uint) error
	Archive(ctx context.Context, id uint, userID uint) error
	Pin(ctx context.Context, id uint, pinned bool, userID uint) error
}

type announcementService struct {
	repo announcementRepo.AnnouncementRepository
}

func NewAnnouncementService(repo announcementRepo.AnnouncementRepository) AnnouncementService {
	return &announcementService{repo: repo}
}

// validateTime 验证发布时间和过期时间
func (s *announcementService) validateTime(publishedAt, expiredAt *time.Time) error {
	if publishedAt != nil && expiredAt != nil {
		if !expiredAt.After(*publishedAt) {
			return apperrors.ErrExpiredTimeInvalid
		}
	}
	return nil
}
