package plugin

import (
	"context"
	"errors"
	"fmt"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// DeletePlugin 删除插件（软删除）
func (s *pluginService) DeletePlugin(ctx context.Context, pluginID, userID uint) error {
	// 1. 查询插件是否存在（包括已软删除的，不允许重复删除）
	plugin, err := s.repo.GetByID(ctx, pluginID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("query plugin failed: %w", err)
	}

	// 3. 如果插件已启用，不允许删除（或需要先停用）
	if plugin.Enabled {
		return apperrors.ErrPluginEnabledFirst
	}

	// 4. 逻辑删除
	if err := s.repo.DeleteByID(ctx, pluginID); err != nil {
		return fmt.Errorf("delete plugin failed: %w", err)
	}

	return nil
}
