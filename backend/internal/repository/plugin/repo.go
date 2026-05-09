package plugin

import (
	"context"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type PluginRepository interface {
	List(ctx context.Context, queryDO *common.PageQuery[do.PluginMeta]) ([]*do.PluginMeta, int64, error) // 分页查询
	Create(ctx context.Context, plugin *do.PluginMeta) error                                             // 创建
	GetByName(ctx context.Context, name string) (*do.PluginMeta, error)                                  // 获取插件
	ListByAuthorID(ctx context.Context, authorID uint) ([]*do.PluginMeta, error)                         // 获取作者的所有插件
	ListEnabled(ctx context.Context) ([]*do.PluginMeta, error)                                           // 获取所有已启用的插件
	Update(ctx context.Context, plugin *do.PluginMeta) error                                             // 更新
	DeleteByID(ctx context.Context, id uint) error                                                       // 硬删除
	IsExist(ctx context.Context, name string) (bool, error)                                              // 判断插件是否存在
	GetByID(ctx context.Context, id uint) (*do.PluginMeta, error)                                        // 软删除后查询
	FindByNameUnscoped(ctx context.Context, name string) (*do.PluginMeta, error)                         // 软删除后查询
	DeletePermanentlyByID(ctx context.Context, id uint) error                                            // 彻底删除
}

type pluginRepo struct {
	db *gorm.DB
}

func NewPluginRepository(db *gorm.DB) PluginRepository {
	return &pluginRepo{db: db}
}

func (r *pluginRepo) FindByNameUnscoped(ctx context.Context, name string) (*do.PluginMeta, error) {
	var plugin do.PluginMeta
	err := r.db.WithContext(ctx).Unscoped().
		Where("name = ?", name).
		First(&plugin).Error
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}
