package attachment

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/vo"
)

func (s *service) ListUserPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[vo.PluginMetaVO], error) {

	return s.pluginSvc.ListPlugins(ctx, queryBO)
}
