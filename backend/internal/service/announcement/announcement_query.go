package announcement

import (
	"context"
	"errors"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

func (s *announcementService) GetByID(ctx context.Context, id uint) (*do.Announcement, error) {
	announcement, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrAnnouncementNotFound
		}
		return nil, err
	}
	go s.repo.IncrementViewCount(context.Background(), id)
	return announcement, nil
}

func (s *announcementService) List(ctx context.Context, req *request.ListAnnouncements) (*vo.ListAnnouncements, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	repoReq := &request.ListAnnouncements{

		Page:      req.Page,
		PageSize:  req.PageSize,
		BoardID:   req.BoardID,
		Type:      req.Type,
		Status:    req.Status,
		IsPinned:  req.IsPinned,
		IsGlobal:  req.IsGlobal,
		Keyword:   req.Keyword,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}
	announcements, total, err := s.repo.List(ctx, repoReq)
	if err != nil {
		return nil, err
	}
	return &vo.ListAnnouncements{
		Total:         total,
		Page:          req.Page,
		PageSize:      req.PageSize,
		Announcements: announcements,
	}, nil
}

func (s *announcementService) GetPinned(ctx context.Context, boardID *uint) ([]do.Announcement, error) {
	return s.repo.GetPinned(ctx, boardID)
}
