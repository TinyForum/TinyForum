package request

import "tiny-forum/internal/model/do"

type ListPluginsRequest struct {
	PageRequest
	Keyword  string            `form:"keyword"`                 // 关键字
	Status   string            `form:"status" default:"active"` // 插件状态
	Tags     []string          `form:"tag"`                     // 标签名称
	AuthorID uint              `form:"author_id"`               // 用户ID
	Type     string            `form:"type"`                    // 插件类型
	Version  string            `form:"version"`                 // 插件版本
	Category do.PluginCategory `json:"category"`                // 插件分类
}

type PluginListRequest struct {
	Page     int             `json:"page" form:"page" default:"1"`
	PageSize int             `json:"page_size" form:"page_size" binding:"min=1,max=100"  default:"20"`
	AuthorID uint            `json:"author_id" form:"author_id"`
	Tags     []string        `json:"tags" form:"tags"`
	Type     string          `json:"type" form:"type"`
	Keyword  string          `json:"keyword" form:"keyword"`
	SortBy   string          `json:"sort_by" form:"sort_by" default:"id"`
	Status   do.PluginStatus `json:"status" form:"status" default:"active"`
}

// UpdatePluginRequest 更新插件元数据的请求参数
type UpdatePluginRequest struct {
	Name        string                 `json:"name" binding:"omitempty,max=100"`        // 插件名称
	Version     string                 `json:"version" binding:"omitempty,max=50"`      // 版本号
	Description string                 `json:"description" binding:"omitempty,max=2000"` // 描述
	Summary     string                 `json:"summary" binding:"omitempty,max=500"`     // 摘要
	Enabled     *bool                  `json:"enabled"`                                 // 是否启用
	Slots       []string               `json:"slots"`                                   // 注入插槽
	Config      map[string]interface{} `json:"config"`                                  // 插件配置
}
