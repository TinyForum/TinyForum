package reports

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/repository/reports"
)

type ReportsService interface {
	Lists(ctx context.Context, listPostsBO *common.PageQuery[bo.ListReportBO]) ([]do.Report, int64, error)
}

type reportsService struct {
	reportRepo reports.ReportsRepository
}

func NewReportsService(
	reportsRepo reports.ReportsRepository,
) ReportsService {
	return &reportsService{
		reportRepo: reportsRepo,
	}
}

func (s *reportsService) Lists(ctx context.Context, listPostsBO *common.PageQuery[bo.ListReportBO]) ([]do.Report, int64, error) {
	return s.reportRepo.List(ctx, listPostsBO)
}
