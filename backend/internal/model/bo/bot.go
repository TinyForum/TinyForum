package bo

import "tiny-forum/internal/model/do"

type CreateBotBO struct {
	// 基础信息
	Name        string   `json:"name" binding:"required,min=1,max=100"`         // 机器人名称
	Version     string   `json:"version" binding:"required,min=1,max=50"`       // 版本号（如 "1.0.0"）
	Description string   `json:"description,omitempty" binding:"max=10000"`     // 详细描述
	Summary     string   `json:"summary,omitempty" binding:"max=300"`           // 一句话简介
	AvatarURL   string   `json:"avatarUrl,omitempty" binding:"omitempty,url"`   // 头像URL
	Screenshots []string `json:"screenshots,omitempty"`                         // 截图URL列表
	HomepageURL string   `json:"homepageUrl,omitempty" binding:"omitempty,url"` // 项目/官网地址

	// 类型与标签
	Type do.BotType `json:"type" binding:"required,oneof=chat moderate notify sync task webhook analysis game"` // 机器人类型
	Tags []string   `json:"tags,omitempty" binding:"max=10,dive,min=1,max=30"`                                  // 最多10个标签，每个标签1-30字符

	// 运行配置
	ScriptURL   string            `json:"scriptUrl" binding:"required,url"`                                   // 机器人代码入口URL
	TriggerType do.BotTriggerType `json:"triggerType" binding:"required,oneof=schedule event webhook manual"` // 触发方式
	CronExpr    string            `json:"cronExpr,omitempty"`                                                 // 定时表达式（当triggerType=schedule时必填）
	EventFilter string            `json:"eventFilter,omitempty"`                                              // 事件过滤条件（JSON字符串，当triggerType=event时使用）

	// 定价与商业化（可选，若不提供则默认为免费）
	Pricing *do.BotPricing `json:"pricing,omitempty"`

	// 权限申请
	Permissions []do.BotPermission `json:"permissions,omitempty" binding:"max=20,dive,oneof=read:user read:posts write:posts read:comments write:comments send:message manage:content read:stats"`

	// 配置模式（机器人开发者定义的用户可配置项，仅首次创建时需要定义schema）
	ConfigSchema []do.BotConfigField `json:"configSchema,omitempty" binding:"max=50"`

	// 初始配置值（机器人创建者自己的配置，按configSchema填入值）
	ConfigValues map[string]any `json:"configValues,omitempty"`
}
