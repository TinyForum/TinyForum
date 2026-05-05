package plugin

import (
	"context"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
)

func (r *pluginRepository) List(ctx context.Context, query *dto.PluginQueryDTO, pageParam common.PageParam) (*common.PageResult[do.PluginMeta], error) {
	db := r.db.WithContext(ctx).Model(&do.PluginMeta{})

	// 动态条件
	if query.AuthorID != 0 {
		db = db.Where("author_id = ?", query.AuthorID)
	}
	if query.Tags != nil {
		db = db.Where("tag_id = ?", query.Tags)
	}
	if query.Type != "" {
		db = db.Where("post_type = ?", query.Type)
	}
	if query.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+query.Keyword+"%")
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// 总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// 列表
	var list []do.PluginMeta
	offset := (pageParam.Page - 1) * pageParam.PageSize
	order := "created_at DESC"
	if query.SortBy != "" {
		order = query.SortBy
	}
	err := db.Offset(offset).Limit(pageParam.PageSize).Order(order).Find(&list).Error
	if err != nil {
		return nil, err
	}

	hasMore := int64(pageParam.Page*pageParam.PageSize) < total

	return &common.PageResult[do.PluginMeta]{
		Total:    total,
		Page:     pageParam.Page,
		PageSize: pageParam.PageSize,
		List:     list,
		HasMore:  hasMore,
	}, nil
}
