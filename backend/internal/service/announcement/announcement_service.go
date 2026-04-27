package announcement

import (
	"context"
	"time"
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	announcementRepo "tiny-forum/internal/repository/announcement"
	apperrors "tiny-forum/pkg/errors"
)

type AnnouncementService interface {
	Create(ctx context.Context, req *dto.CreateAnnouncementRequest, userID uint) (*model.Announcement, error)
	Update(ctx context.Context, id uint, req *dto.UpdateAnnouncementRequest, userID uint) error
	Delete(ctx context.Context, id uint, userID uint) error
	GetByID(ctx context.Context, id uint) (*model.Announcement, error)
	List(ctx context.Context, req *dto.ListAnnouncementRequest) (*dto.ListAnnouncementResponse, error)
	GetPinned(ctx context.Context, boardID *uint) ([]model.Announcement, error)
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
