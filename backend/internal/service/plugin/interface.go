// package plugin

// import (
// 	"context"
// 	"mime/multipart"
// 	"tiny-forum/internal/model/bo"
// 	"tiny-forum/internal/model/common"
// 	"tiny-forum/internal/model/do"
// 	"tiny-forum/internal/model/vo"
// 	pluginRepo "tiny-forum/internal/repository/plugin"
// 	"tiny-forum/internal/storage"
// )

// type pluginService struct {
// 	repo pluginRepo.PluginRepository
// 	 storage storage.StorageDriver
// }

// type PluginService interface {
// 	ListPlugins(ctx context.Context, queryBO *bo.PluginQueryBO) (*common.PageResult[vo.PluginMetaVO], error)
// 	// // 安装插件：上传 zip 包，解压，验证 manifest，写入 db
//     Install(ctx context.Context, fileHeader *multipart.FileHeader, userID uint) (*do.PluginMeta, error)
//     // // 列出所有启用的系统插件
//     // ListAllEnabled(ctx context.Context) ([]*do.PluginMeta, error)
//     // // 列出用户已安装的插件（关联表，如果没有单独的表，可以从 plugin 表里通过 author_id 筛选）
//     ListUserPlugins(ctx context.Context, userID int64) ([]*do.PluginMeta, error)
//     // // 启用/禁用插件
//     // SetEnabled(ctx context.Context, pluginID string, enabled bool) error
// }

//	func NewPluginService(repo pluginRepo.PluginRepository) PluginService {
//		return &pluginService{repo: repo}
//	}
package plugin