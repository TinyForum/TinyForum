package board

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
)

// Count 统计板块总数
func (r *boardRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&do.Board{}).Count(&count).Error
	return count, err
}

// CountByDateRange 按日期范围统计新增板块数
func (r *boardRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&do.Board{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}

// GetHotBoardsByDateRange 获取热门板块（委托给 stats 仓库）
func (r *boardRepository) GetHotBoardsByDateRange(
	ctx context.Context,
	startDate, endDate time.Time,
	limit int,
) ([]*vo.HotBoardRowVO, error) {
	return r.stats.GetHotBoardsByDateRange(ctx, startDate, endDate, limit)
}
