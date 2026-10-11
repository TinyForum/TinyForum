package reports

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	"tiny-forum/internal/repository/reports"
)

type ReportsService interface {
	Lists(ctx context.Context, listPostsBO *common.PageQuery[bo.ListReportBO]) ([]do.Report, int64, error)
	// CreateReport 提交一条举报
	CreateReport(ctx context.Context, reporterID uint, req *request.CreateReportRequest) (*vo.ReportVO, error)
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
