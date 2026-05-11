package plugin

import (
	"context"
	"tiny-forum/internal/model/do"
)

// 彻底删除
func (r *pluginRepo) DeletePermanentlyByID(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&do.PluginMeta{}, id).Error
}

// 逻辑删除（若模型有软删除）或物理删除（若无软删除）
func (r *pluginRepo) DeleteByID(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&do.PluginMeta{}).Error
}
