package plugin

import (
	"context"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type PluginRepository interface {
	Create(ctx context.Context, plugin *do.PluginMeta) error
	GetByName(ctx context.Context, name string) (*do.PluginMeta, error)
	ListByAuthorID(ctx context.Context, authorID int64) ([]*do.PluginMeta, error)
	ListEnabled(ctx context.Context) ([]*do.PluginMeta, error)
	Update(ctx context.Context, plugin *do.PluginMeta) error
	Delete(ctx context.Context, name string) error
}

type pluginRepo struct {
	db *gorm.DB
}

func NewPluginRepository(db *gorm.DB) PluginRepository {
	return &pluginRepo{db: db}
}

func (r *pluginRepo) Create(ctx context.Context, plugin *do.PluginMeta) error {
	return r.db.WithContext(ctx).Create(plugin).Error
}

func (r *pluginRepo) GetByName(ctx context.Context, name string) (*do.PluginMeta, error) {
	var p do.PluginMeta
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&p).Error
	return &p, err
}

func (r *pluginRepo) ListByAuthorID(ctx context.Context, authorID int64) ([]*do.PluginMeta, error) {
	var list []*do.PluginMeta
	err := r.db.WithContext(ctx).Where("author_id = ?", authorID).Find(&list).Error
	return list, err
}

func (r *pluginRepo) ListEnabled(ctx context.Context) ([]*do.PluginMeta, error) {
	var list []*do.PluginMeta
	err := r.db.WithContext(ctx).Where("enabled = ?", true).Find(&list).Error
	return list, err
}

func (r *pluginRepo) Update(ctx context.Context, plugin *do.PluginMeta) error {
	return r.db.WithContext(ctx).Save(plugin).Error
}

func (r *pluginRepo) Delete(ctx context.Context, name string) error {
	return r.db.WithContext(ctx).Where("name = ?", name).Delete(&do.PluginMeta{}).Error
}