package bot

import (
	"context"
	"tiny-forum/internal/model/do"
)

// ListByUser 获取用户创建的机器人
func (r *repo) ListByUser(ctx context.Context, creatorID uint, offset, limit int) ([]*do.Bot, int64, error) {
	var bots []*do.Bot
	var total int64
	query := r.db.WithContext(ctx).Model(&do.Bot{})
	if creatorID > 0 {
		query = query.Where("creator_id = ?", creatorID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&bots).Error
	return bots, total, err
}

// List 获取所有机器人
func (r *repo) List(ctx context.Context, offset, limit int) ([]*do.Bot, int64, error) {
	var bots []*do.Bot
	var total int64
	query := r.db.WithContext(ctx).Model(&do.Bot{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&bots).Error
	return bots, total, err
}

// ListActive 获取所有启用的机器人
func (r *repo) ListActive(ctx context.Context) ([]*do.Bot, error) {
	var bots []*do.Bot
	err := r.db.WithContext(ctx).Where("enabled = ?", true).Find(&bots).Error
	return bots, err
}
