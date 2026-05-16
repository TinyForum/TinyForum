package plugin

import (
	"context"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *pluginRepo) Update(ctx context.Context, plugin *do.PluginManifest) error {
	return r.db.WithContext(ctx).Save(plugin).Error
}

func (r *pluginRepo) TogglePluginStatus(ctx context.Context, pluginSlug string) error {
	// 直接使用 SQL 表达式取反
	return r.db.WithContext(ctx).
		Model(&do.PluginManifest{}).
		Where("slug = ?", pluginSlug).
		Update("enabled", gorm.Expr("NOT enabled")).
		Error
}
