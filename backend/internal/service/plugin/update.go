package plugin

import "context"

func (s *pluginService) TogglePluginStatus(ctx context.Context, pluginID uint) error {
	return s.repo.TogglePluginStatus(ctx, pluginID)
}
