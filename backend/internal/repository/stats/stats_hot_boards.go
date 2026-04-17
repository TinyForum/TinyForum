package stats

import (
	"context"
	"time"
)

// HotBoardRow 热门板块查询结果行
type HotBoardRow struct {
	ID           int64
	Name         string
	Icon         string
	ArticleCount int64
	CommentCount int64
	ActiveUser   int64
}

// GetHotBoardsByDateRange 查询指定时间段内热门板块（按活跃度排序）
func (r *StatsRepository) GetHotBoardsByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*HotBoardRow, error) {
	var rows []*HotBoardRow

	err := r.db.WithContext(ctx).
		Table("boards b").
		Select(`
			b.id,
			b.name,
			b.icon,
			COUNT(DISTINCT p.id)        AS article_count,
			COUNT(DISTINCT c.id)        AS comment_count,
			COUNT(DISTINCT p.author_id) AS active_user
		`).
		Joins("LEFT JOIN posts p ON p.board_id = b.id AND p.deleted_at IS NULL AND p.created_at BETWEEN ? AND ?", startDate, endDate).
		Joins("LEFT JOIN comments c ON c.post_id = p.id AND c.deleted_at IS NULL AND c.created_at BETWEEN ? AND ?", startDate, endDate).
		Where("b.deleted_at IS NULL").
		Group("b.id, b.name, b.icon").
		Order("(COUNT(DISTINCT p.id) * 10 + COUNT(DISTINCT c.id) * 2 + COUNT(DISTINCT p.author_id) * 5) DESC").
		Limit(limit).
		Scan(&rows).Error

	return rows, err
}
