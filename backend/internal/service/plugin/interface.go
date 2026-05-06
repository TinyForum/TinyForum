package plugin

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/vo"
	pluginRepo "tiny-forum/internal/repository/plugin"
)

type PluginService interface {
	ListPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[vo.PluginMetaVO], error)
}

type pluginService struct {
	repo pluginRepo.PluginRepository
}

func NewPluginService(repo pluginRepo.PluginRepository) PluginService {
	return &pluginService{repo: repo}
}
