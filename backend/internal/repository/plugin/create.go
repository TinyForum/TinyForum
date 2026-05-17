package plugin

import (
	"context"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/logger"
)

func (r *pluginRepo) Create(ctx context.Context, plugin *do.PluginManifest) error {
	logger.Infof("创建插件: %v", plugin)
	return r.db.WithContext(ctx).Create(plugin).Error
}
