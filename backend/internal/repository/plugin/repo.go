package plugin

import (
	"context"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type PluginRepository interface {
	// create
	Create(ctx context.Context, plugin *do.PluginManifest) error                                                 // 创建
	List(ctx context.Context, queryDO *common.PageQuery[do.PluginManifest]) ([]*do.PluginManifest, int64, error) // 分页查询
	IsExist(ctx context.Context, name string) (bool, error)                                                      // 判断插件是否存在
	// query
	GetByName(ctx context.Context, name string) (*do.PluginManifest, error)          // 获取插件
	ListByAuthorID(ctx context.Context, authorID uint) ([]*do.PluginManifest, error) // 获取作者的所有插件
	ListEnabled(ctx context.Context) ([]*do.PluginManifest, error)                   // 获取所有已启用的插件
	GetBySlug(ctx context.Context, slug string) (*do.PluginManifest, error)
	// update
	Update(ctx context.Context, plugin *do.PluginManifest) error     // 更新
	TogglePluginStatus(ctx context.Context, pluginSlug string) error // 切换插件状态
	// delete
	DeleteByID(ctx context.Context, id uint) error                                   // 硬删除
	GetByID(ctx context.Context, id uint) (*do.PluginManifest, error)                // 软删除后查询
	FindByNameUnscoped(ctx context.Context, name string) (*do.PluginManifest, error) // 软删除后查询
	DeletePermanentlyByID(ctx context.Context, id uint) error                        // 彻底删除
}

type pluginRepo struct {
	db *gorm.DB
}

func NewPluginRepository(db *gorm.DB) PluginRepository {
	return &pluginRepo{db: db}
}

func (r *pluginRepo) FindByNameUnscoped(ctx context.Context, name string) (*do.PluginManifest, error) {
	var plugin do.PluginManifest
	err := r.db.WithContext(ctx).Unscoped().
		Where("name = ?", name).
		First(&plugin).Error
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}
