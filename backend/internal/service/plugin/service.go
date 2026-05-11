package plugin

import (
	"context"
	"mime/multipart"

	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
	pluginRepo "tiny-forum/internal/repository/plugin"
	"tiny-forum/internal/storage"
)

type pluginService struct {
	repo    pluginRepo.PluginRepository
	storage storage.StorageDriver
	cfg     *config.ConfigPlugins
}
type PluginService interface {
	// create
	Create(ctx context.Context, fileHeader *multipart.FileHeader, userID uint) (*do.PluginMeta, error) // 创建插件（用户上传）
	// query
	ListPlugins(ctx context.Context, queryBO *bo.PageQuery[bo.PluginQueryBO]) (*common.PageResult[vo.PluginMetaVO], error) // 获取插件列表
	ListUserPlugins(ctx context.Context, userID uint) ([]*do.PluginMeta, error)                                            // 获取当前用户创建的插件列表
	// delete
	DeletePlugin(ctx context.Context, pluginID uint, userID uint) error // 删除插件
	// update
	TogglePluginStatus(ctx context.Context, pluginID uint) error // 切换插件状态
}

func NewPluginService(repo pluginRepo.PluginRepository, storage storage.StorageDriver, cfg *config.ConfigPlugins) PluginService {
	return &pluginService{repo: repo, storage: storage, cfg: cfg}
}
