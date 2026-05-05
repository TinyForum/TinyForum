// converter/plugin_converter.go
package converter

import (
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/vo"
)

// ========== BO ↔ QueryDO ==========
func PluginListBOToQueryDTO(b *bo.PluginQueryBO) *dto.PluginQueryDTO {
	if b == nil {
		return nil
	}
	return &dto.PluginQueryDTO{
		AuthorID: b.AuthorID,
		Tags:     b.Tags,
		Type:     b.Type,
		Keyword:  b.Keyword,
		SortBy:   b.SortBy,
		Status:   b.Status,
	}
}

// ========== DO ↔ BO ==========
func PluginDOToBO(d *do.PluginMeta) *bo.PluginMeta {
	if d == nil {
		return nil
	}
	return &bo.PluginMeta{
		ID:        d.ID,
		Name:      d.Name,
		AuthorID:  d.AuthorID,
		Tags:      d.Tags,
		Type:      d.Type,
		Status:    d.Status,
		CreatedAt: d.CreatedAt,
		// 注意：Keyword、SortBy 等查询字段不包含在 DO 中，不需要赋值
	}
}

func PluginBOToDO(b *bo.PluginMeta) *do.PluginMeta {
	if b == nil {
		return nil
	}
	return &do.PluginMeta{
		BaseModel: common.BaseModel{ // 通过嵌入类型名初始化
			ID:        b.ID,
			CreatedAt: b.CreatedAt,
			// UpdatedAt 和 DeletedAt 通常由数据库自动维护，不需要从 BO 传递
		},
		Name:     b.Name,
		AuthorID: b.AuthorID,
		Tags:     b.Tags,
		Type:     b.Type,
		Status:   b.Status,
	}
}

// ========== PageResult[DO] ↔ PageResult[BO] ==========
func PageDOToPageBO(pageDO *common.PageResult[do.PluginMeta], mapper func(*do.PluginMeta) *bo.PluginMeta) *common.PageResult[bo.PluginMeta] {
	if pageDO == nil {
		return nil
	}
	listBO := make([]bo.PluginMeta, 0, len(pageDO.List))
	for _, d := range pageDO.List {
		// 注意：d 是值类型，取地址传给 mapper
		if boItem := mapper(&d); boItem != nil {
			listBO = append(listBO, *boItem)
		}
	}
	return &common.PageResult[bo.PluginMeta]{
		Total:    pageDO.Total,
		Page:     pageDO.Page,
		PageSize: pageDO.PageSize,
		List:     listBO,
		HasMore:  pageDO.HasMore,
	}
}

// ========== BO ↔ VO ==========
func PluginBOToVO(b *bo.PluginMeta) *vo.PluginVO {
	if b == nil {
		return nil
	}
	return &vo.PluginVO{
		ID:        b.ID,
		Name:      b.Name,
		AuthorID:  b.AuthorID,
		CreatedAt: b.CreatedAt,
		Status:    b.Status,
	}
}

// ========== Request ↔ BO ==========
func PluginListRequestToBO(req *vo.PluginListRequest) *bo.PluginQueryBO {
	if req == nil {
		return nil
	}
	return &bo.PluginQueryBO{
		Page:     req.Page,
		PageSize: req.PageSize,
		AuthorID: req.AuthorID,
		Tags:     req.Tags,
		Type:     req.PostType,
		Keyword:  req.Keyword,
		SortBy:   req.SortBy,
		Status:   req.Status,
	}
}

// 泛型辅助：转换整个 PageResult[BO] -> PageResult[VO]
func PageBOToPageVO(pageBO *common.PageResult[bo.PluginMeta], mapper func(*bo.PluginMeta) *vo.PluginVO) *common.PageResult[vo.PluginVO] {
	if pageBO == nil {
		return nil
	}
	listVO := make([]vo.PluginVO, 0, len(pageBO.List))
	for _, b := range pageBO.List {
		// b 是值类型，取地址传入 mapper
		if voItem := mapper(&b); voItem != nil {
			listVO = append(listVO, *voItem)
		}
	}
	return &common.PageResult[vo.PluginVO]{
		Total:    pageBO.Total,
		Page:     pageBO.Page,
		PageSize: pageBO.PageSize,
		List:     listVO,
		HasMore:  pageBO.HasMore,
	}
}

func PluginQueryBOToQueryDO(b *bo.PluginQueryBO) *dto.PluginQueryDTO {
	if b == nil {
		return nil
	}
	return &dto.PluginQueryDTO{
		Name:     b.Name,
		Type:     b.Type,
		Category: b.Category,
		Tags:     b.Tags,
		AuthorID: b.AuthorID,
		Status:   b.Status,
		// Enabled:  b.Enabled,
		Keyword: b.Keyword,
	}
}
