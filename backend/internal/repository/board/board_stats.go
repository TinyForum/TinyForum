package board

import (
	"context"
	"time"
	"tiny-forum/internal/model"
	statsRepo "tiny-forum/internal/repository/stats"
)

// Count 统计板块总数
func (r *BoardRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Board{}).Count(&count).Error
	return count, err
}

// CountByDateRange 按日期范围统计新增板块数
func (r *BoardRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Board{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

// GetHotBoardsByDateRange 获取热门板块（委托给 stats 仓库）
func (r *BoardRepository) GetHotBoardsByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*statsRepo.HotBoardRow, error) {
	return r.stats.GetHotBoardsByDateRange(ctx, startDate, endDate, limit)
}
