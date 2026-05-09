// package plugin

// import (
// 	"context"
// 	"tiny-forum/internal/model/bo"
// 	"tiny-forum/internal/model/common"
// 	"tiny-forum/internal/model/converter"
// 	"tiny-forum/internal/model/do"
// 	"tiny-forum/internal/model/vo"
// )

// // service/plugin_service.go
// func (s *pluginService) ListPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[vo.PluginMetaVO], error) {
// 	// BO -> Query DO
// 	queryDO := converter.PluginQueryBOToQueryDO(queryBO)

// 	// 调用 Repo
// 	pageDO, err := s.repo.List(ctx, queryDO, common.PageParam{
// 		Page:     queryBO.Page,
// 		PageSize: queryBO.PageSize,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// DO Page -> BO Page
// 	pageVO := converter.PageDOToPageVO(pageDO)
// 	return pageVO, nil
// }

// // ListUserPlugins 获取用户已安装的插件（通过 author_id）
//
//	func (s *pluginService) ListUserPlugins(ctx context.Context, userID int64) ([]*do.PluginMeta, error) {
//	    return s.repo.ListByAuthorID(ctx, userID)
//	}
package plugin