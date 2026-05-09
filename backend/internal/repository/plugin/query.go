package plugin

import (
	"context"
	"errors"
	"fmt"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/logger"

	"gorm.io/gorm"
)

// GetByName 根据名称获取未软删除的插件
func (r *pluginRepo) GetByName(ctx context.Context, name string) (*do.PluginMeta, error) {
	var p do.PluginMeta
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&p).Error
	return &p, err
}

// IsExist 检查插件是否存在
func (r *pluginRepo) IsExist(ctx context.Context, name string) (bool, error) {
	var p do.PluginMeta
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&p).Error
	if err == nil {
		return true, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	logger.Errorf("check plugin existence failed: %v", err)
	return false, err
}

// Create 创建插件
func (r *pluginRepo) ListByAuthorID(ctx context.Context, authorID uint) ([]*do.PluginMeta, error) {
	var list []*do.PluginMeta
	err := r.db.WithContext(ctx).Where("author_id = ?", authorID).Find(&list).Error
	return list, err
}

// ListEnabled 获取所有已启用的插件
func (r *pluginRepo) ListEnabled(ctx context.Context) ([]*do.PluginMeta, error) {
	var list []*do.PluginMeta
	err := r.db.WithContext(ctx).Where("enabled = ?", true).Find(&list).Error
	return list, err
}

// GetByID 根据ID获取未软删除的插件
func (r *pluginRepo) GetByID(ctx context.Context, id uint) (*do.PluginMeta, error) {
	var plugin do.PluginMeta
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&plugin).Error
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}

func (r *pluginRepo) List(ctx context.Context, queryBO *common.PageQuery[do.PluginMeta]) ([]*do.PluginMeta, int64, error) {
	db := r.db.WithContext(ctx).Model(&do.PluginMeta{})

	// 软删除过滤
	db = db.Where("deleted_at IS NULL")

	// 用户过滤
	if queryBO.Data.AuthorID > 0 {
		db = db.Where("user_id = ?", queryBO.Data.AuthorID)
	}

	// 状态过滤
	if queryBO.Data.Status != "" {
		db = db.Where("status = ?", queryBO.Data.Status)
	}

	// 标签过滤
	if len(queryBO.Data.Tags) != 0 {
		db = db.Where("tag = ?", queryBO.Data.Tags)
	}

	// 关键字模糊搜索
	if queryBO.Keyword != "" {
		keyword := "%" + queryBO.Keyword + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	// 排序
	sortField := queryBO.SortBy
	if sortField == "" {
		sortField = "created_at"
	}
	// 默认降序
	order := queryBO.Order
	if order == "" {
		order = "desc"
	}
	// 支持自定义排序字段
	db = db.Order(fmt.Sprintf("%s %s", sortField, order))

	// 分页查询
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err

	}

	var plugins []*do.PluginMeta
	offset := (queryBO.Page - 1) * queryBO.PageSize
	err := db.Offset(offset).Limit(queryBO.PageSize).Find(&plugins).Error
	return plugins, total, err

}
