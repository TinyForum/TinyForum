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
