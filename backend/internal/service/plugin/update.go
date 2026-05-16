package plugin

import "context"

func (s *pluginService) TogglePluginStatus(ctx context.Context, pluginSlug string) error {
	return s.repo.TogglePluginStatus(ctx, pluginSlug)
}
