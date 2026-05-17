package admin

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

func (s *adminService) ListApplications(boardID *uint, status do.ApplicationStatus, page, pageSize int) ([]do.ModeratorApplication, int64, error) {
	return s.boardSvc.ListApplications(boardID, status, page, pageSize)
}

func (s *adminService) ReviewApplication(ctx context.Context, input request.ReviewApplicationRequest, reviewerID uint) error {
	return s.boardSvc.ReviewApplication(ctx, input, reviewerID)
}
