package admin

import (
	"context"
	"tiny-forum/internal/model/query"
	"tiny-forum/internal/model/vo"
)

func (s *adminService) ListAnnouncements(ctx context.Context, req *query.ListAnnouncements) (*vo.ListAnnouncements, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	repoReq := &query.ListAnnouncements{
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
	announcements, total, err := s.announcementRepo.List(ctx, repoReq)
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
