package tag

import (
	"context"
	"time"
	"tiny-forum/internal/model"
)

// Count 获取标签总数
func (r *tagRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Model(&model.Tag{}).Count(&count).Error
	return count, err
}

// CountByDateRange 根据日期范围统计标签数
func (r *tagRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&model.Tag{}).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Count(&count).Error
	return count, err
}
