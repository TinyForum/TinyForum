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

// PluginMetaVO 插件元数据脱敏视图（对外暴露）
type PluginMetaVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 基础标识
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	IconURL     string   `json:"iconUrl,omitempty"`
	Screenshots []string `json:"screenshots,omitempty"`
	HomepageURL string   `json:"homepageUrl,omitempty"`

	// 分类与类型
	Type     string   `json:"type"`     // PluginType 映射为字符串
	Category string   `json:"category"` // PluginCategory 映射为字符串
	Tags     []string `json:"tags,omitempty"`

	// 作者信息（移除邮箱）
	AuthorID  uint   `json:"authorId,omitempty"`
	AuthorURL string `json:"authorUrl,omitempty"`

	// 加载配置
	ScriptURL   string   `json:"scriptUrl"`
	ServerEntry string   `json:"serverEntry,omitempty"`
	Slots       []string `json:"slots,omitempty"`
	Routes      []string `json:"routes,omitempty"`

	// 价格与兼容性
	Pricing       interface{} `json:"pricing,omitempty"` // 使用 interface{} 或自定义 VO
	Compatibility interface{} `json:"compatibility,omitempty"`

	// 权限
	Permissions []interface{} `json:"permissions,omitempty"` // 权限声明列表

	// 运行时
	Enabled      bool    `json:"enabled"`
	Status       string  `json:"status"` // PluginStatus 映射为字符串
	InstallCount int     `json:"installCount"`
	Rating       float32 `json:"rating"`

	// 配置（仅保留 Schema，不返回具体配置值）
	ConfigSchema []interface{} `json:"configSchema,omitempty"`
}
