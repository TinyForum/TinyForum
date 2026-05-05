package dto

import (
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
)

// PluginMeta 插件元数据（数据库模型）
type PluginMeta struct {
	// 基础信息
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // 插入时自动填充
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// 基础标识
	Name        string   `json:"name" gorm:"type:varchar(100);not null;index:idx_name,unique"` // 插件名称
	Version     string   `json:"version" gorm:"type:varchar(50);not null"`                     // 插件版本
	Description string   `json:"description" gorm:"type:text"`                                 // 插件描述
	Summary     string   `json:"summary,omitempty" gorm:"type:varchar(300)"`                   // 一句话简介
	IconURL     string   `json:"iconUrl,omitempty" gorm:"type:varchar(255)"`                   // 插件图标
	Screenshots []string `json:"screenshots,omitempty" gorm:"type:json"`                       // 截图列表
	HomepageURL string   `json:"homepageUrl,omitempty" gorm:"type:varchar(255)"`               // 官网地址

	// 分类与类型
	Type     do.PluginType     `json:"type" gorm:"type:varchar(20);not null;index"`     // 插件类型（对于服务端）
	Category do.PluginCategory `json:"category" gorm:"type:varchar(30);not null;index"` // 插件分类（对于业务）
	Tags     []string          `json:"tags,omitempty" gorm:"type:json"`                 // 标签列表

	// 作者信息
	Author      string `json:"author" gorm:"type:varchar(100);not null"`       // 作者名称
	AuthorEmail string `json:"authorEmail,omitempty" gorm:"type:varchar(100)"` // 作者邮箱
	AuthorURL   string `json:"authorUrl,omitempty" gorm:"type:varchar(255)"`   // 作者主页

	// 加载配置
	ScriptURL   string   `json:"scriptUrl" gorm:"type:varchar(500);not null"`    // 前端脚本入口
	ServerEntry string   `json:"serverEntry,omitempty" gorm:"type:varchar(255)"` // 后端服务入口
	Slots       []string `json:"slots,omitempty" gorm:"type:json"`               // 注入的插槽名称列表
	Routes      []string `json:"routes,omitempty" gorm:"type:json"`              // 注册的路由路径

	// 价格与兼容性
	Pricing       do.PluginPricing       `json:"pricing" gorm:"type:json;serializer:json"`       // 价格信息
	Compatibility do.PluginCompatibility `json:"compatibility" gorm:"type:json;serializer:json"` // 兼容性信息

	// 权限
	Permissions []do.PluginPermission `json:"permissions,omitempty" gorm:"type:json"` // 权限声明

	// 运行时（服务端写入，前端只读）
	Enabled      bool            `json:"enabled" gorm:"default:false;index"`                // 是否启用
	Status       do.PluginStatus `json:"status" gorm:"type:varchar(20);default:'inactive'"` // 插件状态
	InstallCount int             `json:"installCount" gorm:"default:0"`                     // 安装次数
	Rating       float32         `json:"rating" gorm:"type:decimal(2,1);default:0"`         // 评分 0~5

	// 配置
	ConfigSchema []do.PluginConfigField `json:"configSchema,omitempty" gorm:"type:json;serializer:json"` // 配置字段定义
	Config       map[string]any         `json:"config,omitempty" gorm:"type:json;serializer:json"`       // 配置值
}

type PluginList struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// 基础标识
	Name    string `json:"name" gorm:"type:varchar(100);not null;index:idx_name,unique"` // 插件名称
	Version string `json:"version" gorm:"type:varchar(50);not null"`                     // 插件版本
	Summary string `json:"summary,omitempty" gorm:"type:varchar(300)"`                   // 一句话简介
	IconURL string `json:"iconUrl,omitempty" gorm:"type:varchar(255)"`                   // 插件图标

	// 分类与类型
	Type     do.PluginType     `json:"type" gorm:"type:varchar(20);not null;index"`     // 插件类型（对于服务端）
	Category do.PluginCategory `json:"category" gorm:"type:varchar(30);not null;index"` // 插件分类（对于业务）
	Tags     []string          `json:"tags,omitempty" gorm:"type:json"`                 // 标签列表

	// 作者信息
	Author string `json:"author" gorm:"type:varchar(100);not null"` // 作者名称
}

type PluginListOptionsQuery struct {
	request.PageRequest
	AuthorID uint
	TagID    uint
	PostType string
	Keyword  string
	SortBy   string
	Status   do.PluginStatus
}

// type PluginQuery struct {
//     AuthorID uint
//     TagID    uint
//     PostType string
//     Keyword  string
//     SortBy   string
//     Status   do.PluginStatus
// }

type PluginQueryDTO struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Category string          `json:"category"`
	Tags     []string        `json:"tags"`
	AuthorID uint            `json:"author_id"`
	Status   do.PluginStatus `json:"status"`
	SortBy   string          `json:"sort_by"`
	Enabled  *bool           `json:"enabled"` // 使用指针区分未传/传false
	Keyword  string          `json:"keyword"` // 模糊搜索关键词
}
