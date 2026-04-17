package announcement

import (
	"context"
	"time"

	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
	apperrors "tiny-forum/pkg/errors"
)

type AnnouncementService interface {
	Create(ctx context.Context, req *CreateAnnouncementRequest, userID uint) (*model.Announcement, error)
	Update(ctx context.Context, id uint, req *UpdateAnnouncementRequest, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
	GetByID(ctx context.Context, id uint) (*model.Announcement, error)
	List(ctx context.Context, req *ListAnnouncementRequest) (*ListAnnouncementResponse, error)
	GetPinned(ctx context.Context, boardID *uint) ([]model.Announcement, error)
	Publish(ctx context.Context, id uint, userID uint) error
	Archive(ctx context.Context, id uint, userID uint) error
	Pin(ctx context.Context, id uint, pinned bool, userID uint) error
}

type announcementService struct {
	repo repository.AnnouncementRepository
}

func NewAnnouncementService(repo repository.AnnouncementRepository) AnnouncementService {
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
