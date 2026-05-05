package plugin

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/converter"
)

// service/plugin_service.go
func (s *pluginService) ListPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[bo.PluginMeta], error) {
	// BO -> Query DO
	queryDO := converter.PluginQueryBOToQueryDO(queryBO)

	// 调用 Repo
	pageDO, err := s.repo.List(ctx, queryDO, common.PageParam{
		Page:     queryBO.Page,
		PageSize: queryBO.PageSize,
	})
	if err != nil {
		return nil, err
	}

	// DO Page -> BO Page
	pageBO := converter.PageDOToPageBO(pageDO, converter.PluginDOToBO)
	return pageBO, nil
}
