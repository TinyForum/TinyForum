package bo

import (
	"time"
	"tiny-forum/internal/model/do"
)

type PluginMeta struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	AuthorID  uint            `json:"authorId"`
	Tags      []string        `json:"tags"`
	Type      do.PluginType   `json:"type"`
	Keyword   string          `json:"keyword"` // 仅用于查询，不作为返回字段
	SortBy    string          `json:"sortBy"`
	Status    do.PluginStatus `json:"status"`
	CreatedAt time.Time       `json:"createdAt"`
}

// PluginListBO 用于 Service 接收查询参数（不含业务返回字段）
type PluginQueryBO struct {
	Name     string          `json:"name"`
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	AuthorID uint            `json:"authorId"`
	Category string          `json:"category"`
	Tags     []string        `json:"tags"`
	Type     string          `json:"type"`
	Keyword  string          `json:"keyword"`
	SortBy   string          `json:"sortBy"`
	Status   do.PluginStatus `json:"status"`
}

// type PluginQueryBO struct {
//     Page     int      `json:"page"`
//     PageSize int      `json:"pageSize"`
//     Name     string   `json:"name"`
//     Type     string   `json:"type"`
//     Category string   `json:"category"`
//     Tags     []string `json:"tags"`
//     AuthorID uint     `json:"authorId"`
//     Status    do.PluginStatus  `json:"status"`
//     // Enabled  *bool    `json:"enabled"`
//     Keyword  string   `json:"keyword"`
// }
