package article

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
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

func (r *articleRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&do.Article{}).Count(&count).Error
	return count, err
}

func (r *articleRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&do.Article{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

// GetHotArticlesByDateRange 委托给 StatsRepository
func (r *articleRepository) GetHotArticlesByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*vo.HotArticleRowVO, error) {
	return r.stats.GetHotArticlesByDateRange(ctx, startDate, endDate, limit)
}
