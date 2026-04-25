package stats

import (
	"context"
	"time"
	"tiny-forum/internal/model"
)

// HotArticleRow 热门文章查询结果行
type HotArticleRow struct {
	ID           int64
	Title        string
	BoardID      int64
	BoardName    string
	AuthorID     int64
	AuthorName   string
	ViewCount    int64
	CommentCount int64
	LikeCount    int64
}

// GetHotArticlesByDateRange 查询指定时间段内热门文章（按综合热度排序）
func (r *statsRepository) GetHotArticlesByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*HotArticleRow, error) {
	var rows []*HotArticleRow

	err := r.db.WithContext(ctx).
		Table("posts p").
		Select(`
			p.id,
			p.title,
			p.board_id,
			b.name  AS board_name,
			p.author_id,
			u.username AS author_name,
			p.view_count,
			p.like_count,
			COUNT(c.id) AS comment_count
		`).
		Joins("LEFT JOIN boards b ON b.id = p.board_id AND b.deleted_at IS NULL").
		Joins("LEFT JOIN users u ON u.id = p.author_id AND u.deleted_at IS NULL").
		Joins("LEFT JOIN comments c ON c.post_id = p.id AND c.deleted_at IS NULL AND c.created_at BETWEEN ? AND ?", startDate, endDate).
		Where("p.deleted_at IS NULL AND p.status = ? AND p.created_at BETWEEN ? AND ?",
			model.PostStatusPublished, startDate, endDate).
		Group("p.id, b.name, u.username").
		Order("(p.view_count + COUNT(c.id) * 10 + p.like_count * 5) DESC").
		Limit(limit).
		Scan(&rows).Error

	return rows, err
}
