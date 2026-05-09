package plugin

import (
	"context"
	"tiny-forum/internal/model/do"
)

func (r *pluginRepo) Update(ctx context.Context, plugin *do.PluginMeta) error {
	return r.db.WithContext(ctx).Save(plugin).Error
}
