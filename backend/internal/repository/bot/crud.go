package bot

import (
	"context"
	"tiny-forum/internal/model/do"
)

func (r *repo) Create(ctx context.Context, bot *do.Bot) error {
	return r.db.WithContext(ctx).Create(bot).Error
}

func (r *repo) Update(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&do.Bot{}).Where("id = ?", id).Updates(updates).Error
}

func (r *repo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&do.Bot{}, id).Error
}

func (r *repo) GetByID(ctx context.Context, id uint) (*do.Bot, error) {
	var bot do.Bot
	err := r.db.WithContext(ctx).First(&bot, id).Error
	if err != nil {
		return nil, err
	}
	return &bot, nil
}
