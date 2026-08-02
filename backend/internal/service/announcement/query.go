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

func (s *announcementService) List(ctx context.Context, req *request.ListAnnouncementsRequest) (*vo.ListAnnouncements, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	repoReq := &request.ListAnnouncementsRequest{

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

	vos := make([]vo.AnnouncementVO, len(announcements))
	for i, a := range announcements {
		vos[i] = vo.AnnouncementVO{
			ID:          a.ID,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
			Title:       a.Title,
			Content:     a.Content,
			Summary:     a.Summary,
			Cover:       a.Cover,
			Type:        a.Type,
			Status:      a.Status,
			IsPinned:    a.IsPinned,
			IsGlobal:    a.IsGlobal,
			BoardID:     a.BoardID,
			PublishedAt: a.PublishedAt,
			ExpiredAt:   a.ExpiredAt,
			ViewCount:   a.ViewCount,
			CreatedBy:   a.CreatedBy,
		}
		if a.Creator != nil {
			vos[i].Creator.ID = a.Creator.ID
			vos[i].Creator.Username = a.Creator.Username
			vos[i].Creator.AvatarUrl = a.Creator.AvatarUrl
		}
		if a.Board != nil {
			vos[i].Board.ID = a.Board.ID
			vos[i].Board.Name = a.Board.Name
		}
	}

	return &vo.ListAnnouncements{
		Total:         total,
		Page:          req.Page,
		PageSize:      req.PageSize,
		Announcements: vos,
	}, nil
}

func (s *announcementService) GetPinned(ctx context.Context, boardID *uint) ([]do.Announcement, error) {
	return s.repo.GetPinned(ctx, boardID)
}
