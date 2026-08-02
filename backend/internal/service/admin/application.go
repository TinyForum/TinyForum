package admin

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

func (s *adminService) ListApplications(boardID *uint, status do.ApplicationStatus, page, pageSize int) ([]do.ModeratorApplication, int64, error) {
	return s.boardSvc.ListApplications(boardID, status, page, pageSize)
}

func (s *adminService) ReviewApplication(ctx context.Context, input request.ReviewApplicationRequest, reviewerID uint) error {
	return s.boardSvc.ReviewApplication(ctx, input, reviewerID)
}

func (s *adminService) ListReports(ctx context.Context, query *common.PageQuery[bo.ListReportBO]) ([]do.Report, int64, error) {
	return s.reportsSvc.Lists(ctx, query)
}
