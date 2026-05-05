package vo

import (
	"time"
	"tiny-forum/internal/model/do"
)

type ListPlugin struct {
	Total         int64           `json:"total"`
	Page          int             `json:"page"`
	PageSize      int             `json:"page_size"`
	Announcements []do.PluginMeta `json:"announcements"`
}

type PluginListOptions struct {
	AuthorID uint
	Tags     []string
	PostType string
	Keyword  string
	SortBy   string
	Status   do.PluginStatus
}

// 前端返回的视图对象
type PluginVO struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	AuthorID  uint            `json:"authorId"`
	CreatedAt time.Time       `json:"createdAt"`
	Status    do.PluginStatus `json:"status"`
}

// 前端请求参数（也可以放在 request 包，这里遵循你之前的习惯）
type PluginListRequest struct {
	Page     int             `json:"page" form:"page" binding:"min=1"`
	PageSize int             `json:"pageSize" form:"pageSize" binding:"min=1,max=100"`
	AuthorID uint            `json:"authorId" form:"authorId"`
	Tags     []string        `json:"tagId" form:"tags"`
	PostType string          `json:"postType" form:"postType"`
	Keyword  string          `json:"keyword" form:"keyword"`
	SortBy   string          `json:"sortBy" form:"sortBy"`
	Status   do.PluginStatus `json:"status" form:"status"`
}
