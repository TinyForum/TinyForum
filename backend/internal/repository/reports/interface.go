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
	// ListByBoardID 分页查询指定板块下的举报（帖子或其评论属于该板块），可按状态过滤
	ListByBoardID(ctx context.Context, boardID uint, status do.ReportStatus, limit, offset int) ([]do.Report, int64, error)
	// Create 创建一条举报记录
	Create(ctx context.Context, report *do.Report) error
	// GetByID 根据 ID 查询单条举报（预加载举报人与处理人）
	GetByID(ctx context.Context, id uint) (*do.Report, error)
	// Update 更新举报记录（处理举报时写入处理状态）
	Update(ctx context.Context, report *do.Report) error
	// CountIllegalStats 统计指定时间区间内的违规（举报）综合数据
	CountIllegalStats(ctx context.Context, startDate, endDate time.Time) (*IllegalStats, error)
}
type reportsRepository struct {
	db *gorm.DB
}

func NewReportsRepository(db *gorm.DB) ReportsRepository {
	return &reportsRepository{db: db}
}
