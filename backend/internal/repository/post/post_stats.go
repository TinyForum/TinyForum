package post

import (
	"context"
	"time"
	"tiny-forum/internal/model/po"
	"tiny-forum/internal/repository/stats"
)

// HotArticleRow 热门文章行结构（需与 StatsRepository 返回类型匹配）
// type HotArticleRow struct {
// 	ID           uint
// 	Title        string
// 	BoardID      uint
// 	BoardName    string
// 	AuthorID     uint
// 	AuthorName   string
// 	ViewCount    int
// 	CommentCount int
// 	LikeCount    int
// }

func (r *postRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&po.Post{}).Count(&count).Error
	return count, err
}

func (r *postRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&po.Post{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

// GetHotArticlesByDateRange 委托给 StatsRepository
func (r *postRepository) GetHotArticlesByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*stats.HotArticleRow, error) {
	return r.stats.GetHotArticlesByDateRange(ctx, startDate, endDate, limit)
}
