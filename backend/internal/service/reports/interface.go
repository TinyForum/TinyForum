package reports

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	"tiny-forum/internal/repository/reports"
	checkSvc "tiny-forum/internal/service/check"
)

type ReportsService interface {
	Lists(ctx context.Context, listPostsBO *common.PageQuery[bo.ListReportBO]) ([]do.Report, int64, error)
	// CreateReport 提交一条举报
	CreateReport(ctx context.Context, reporterID uint, req *request.CreateReportRequest) (*vo.ReportVO, error)
	// GetByID 获取单条举报详情
	GetByID(ctx context.Context, id uint) (*vo.ReportVO, error)
	// Handle 处理一条举报（通过/驳回）
	Handle(ctx context.Context, id uint, req *request.HandleReportRequest, handlerID uint) (*vo.ReportVO, error)
}

type reportsService struct {
	reportRepo reports.ReportsRepository
	checkSvc   checkSvc.ContentCheckService
}

func NewReportsService(
	reportsRepo reports.ReportsRepository,
	checkSvc checkSvc.ContentCheckService,
) ReportsService {
	return &reportsService{
		reportRepo: reportsRepo,
		checkSvc:   checkSvc,
	}
}
