package plugin

import (
	"context"
	"tiny-forum/internal/model/do"
)

func (r *pluginRepo) Create(ctx context.Context, plugin *do.PluginManifest) error {
	return r.db.WithContext(ctx).Create(plugin).Error
}
