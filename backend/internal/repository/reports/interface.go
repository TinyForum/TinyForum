package reports

import (
	"context"
	"time"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// IllegalStats 违规（举报）聚合统计
type IllegalStats struct {
	Total      int64 // 违规总数
	UserCount  int64 // 涉及用户数
	ArticleCnt int64 // 涉及文章数
	CommentCnt int64 // 涉及评论数
}

type ReportsRepository interface {
	List(ctx context.Context, listReportBO *common.PageQuery[bo.ListReportBO]) ([]do.Report, int64, error)
	// Create 新增一条举报记录
	Create(ctx context.Context, report *do.Report) error
	// CountIllegalStats 统计指定时间段内的违规（举报）聚合数据
	CountIllegalStats(ctx context.Context, startDate, endDate time.Time) (*IllegalStats, error)
}
type reportsRepository struct {
	db *gorm.DB
}

func NewRiskRepository(db *gorm.DB) ReportsRepository {
	return &reportsRepository{db: db}
}
