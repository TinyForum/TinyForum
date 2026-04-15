package repository

import (
	"context"
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// StatsRepository 统计数据仓储
type StatsRepository struct {
	db *gorm.DB
}

func NewStatsRepository(db *gorm.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

// ── 热门文章 ─────────────────────────────────────────────────────────────────

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
func (r *StatsRepository) GetHotArticlesByDateRange(
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

// ── 热门板块 ─────────────────────────────────────────────────────────────────

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
