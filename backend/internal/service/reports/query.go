package reports

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
)

func (s *reportsService) Lists(ctx context.Context, listPostsBO *common.PageQuery[bo.ListReportBO]) ([]do.Report, int64, error) {
	return s.reportRepo.List(ctx, listPostsBO)
}
