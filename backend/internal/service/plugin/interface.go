package plugin

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	pluginRepo "tiny-forum/internal/repository/plugin"
)

type PluginService interface {
	ListPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[bo.PluginMeta], error)
}

type pluginService struct {
	repo pluginRepo.PluginRepository
}

func NewPluginService() PluginService {
	return &pluginService{}
}
