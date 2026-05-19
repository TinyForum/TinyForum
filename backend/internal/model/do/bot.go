package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

type Bot struct {
	common.BaseModel
	Name        string   `json:"name" gorm:"type:varchar(100);not null;index"` // 名称
	Version     string   `json:"version" gorm:"type:varchar(50);not null"`     // 版本
	Description string   `json:"description" gorm:"type:text"`                 // 描述
	Summary     string   `json:"summary" gorm:"type:varchar(300)"`             // 摘要
	AvatarURL   string   `json:"avatar_url" gorm:"type:varchar(255)"`          // 头像
	Screenshots []string `json:"screenshots" gorm:"type:json;serializer:json"` // 截图
	HomepageURL string   `json:"homepage_url" gorm:"type:varchar(255)"`        // 官网

	Type BotType  `json:"type" gorm:"type:varchar(30);not null;index"` // 类型
	Tags []string `json:"tags" gorm:"type:json;serializer:json"`       // 标签

	CreatorID   uint   `json:"creator_id" gorm:"type:bigint;not null;index"` // 创建者 ID
	CreatorName string `json:"creator_name" gorm:"type:varchar(100)"`        // 创建者名称

	ScriptCode string `json:"script_code" gorm:"type:text"`        // 脚本代码
	ScriptURL  string `json:"script_url" gorm:"type:varchar(500)"` // 脚本 URL

	TriggerType BotTriggerType `json:"trigger_type" gorm:"type:varchar(20);not null"` // 触发类型
	CronExpr    string         `json:"cron_expr" gorm:"type:varchar(100)"`            // Cron 表达式
	EventFilter string         `json:"event_filter" gorm:"type:varchar(200)"`         // 事件过滤器

	TimeoutSec    int               `json:"timeout_sec" gorm:"default:10"`                   // 超时时间（秒）
	RetryTimes    int               `json:"retry_times" gorm:"default:0"`                    // 重试次数
	EnvVars       map[string]string `json:"env_vars" gorm:"type:json;serializer:json"`       // 环境变量
	ResourceLimit *ResourceLimit    `json:"resource_limit" gorm:"type:json;serializer:json"` // 资源限制

	Pricing     BotPricing      `json:"pricing" gorm:"type:json;serializer:json"`     // 定价
	Permissions []BotPermission `json:"permissions" gorm:"type:json;serializer:json"` // 权限

	Enabled    bool       `json:"enabled" gorm:"default:false;index"`                // 是否启用
	Status     BotStatus  `json:"status" gorm:"type:varchar(20);default:'inactive'"` // 状态
	ExecCount  int64      `json:"exec_count" gorm:"default:0"`                       // 执行次数
	LastExecAt *time.Time `json:"last_exec_at" gorm:"type:timestamp"`                // 最后执行时间
	ErrorMsg   string     `json:"error_msg" gorm:"type:text"`

	ConfigSchema []BotConfigField `json:"config_schema" gorm:"type:json;serializer:json"` // 配置项
	ConfigValues map[string]any   `json:"config_values" gorm:"type:json;serializer:json"` // 配置值
}

type BotStatus string

const (
	BotStatusActive   BotStatus = "active"   // 激活
	BotStatusInactive BotStatus = "inactive" // 未激活
	BotStatusError    BotStatus = "error"    // 错误
	BotStatusLoading  BotStatus = "loading"  // 加载中
	BotStatusStopped  BotStatus = "stopped"  // 停止
)

// enun [active, inactive, error, loading, stopped]

type BotType string

const (
	BotTypeChat     BotType = "chat"     // 聊天
	BotTypeModerate BotType = "moderate" // 审核
	BotTypeNotify   BotType = "notify"   // 通知
	BotTypeSync     BotType = "sync"     // 同步
	BotTypeTask     BotType = "task"     // 任务
	BotTypeWebhook  BotType = "webhook"  // webhook
	BotTypeAnalysis BotType = "analysis" // 分析
)

//

type BotTriggerType string

const (
	TriggerSchedule BotTriggerType = "schedule" // 定时
	TriggerEvent    BotTriggerType = "event"    // 事件
	TriggerWebhook  BotTriggerType = "webhook"  // webhook
	TriggerManual   BotTriggerType = "manual"   // 手动
)

//

type BotPricingModel string

const (
	BotPricingFree      BotPricingModel = "free"         // 免费
	BotPricingFreemium  BotPricingModel = "freemium"     // 试用
	BotPricingPaid      BotPricingModel = "paid"         // 付费
	PricingSubscription BotPricingModel = "subscription" // 订阅
)

type BotPermission string

const (
	BotPermReadUser      BotPermission = "read:user"      // 读取用户信息
	BotPermReadPosts     BotPermission = "read:posts"     // 读取帖子信息
	BotPermWritePosts    BotPermission = "write:posts"    // 写入帖子信息
	BotPermReadComments  BotPermission = "read:comments"  // 读取评论信息
	BotPermWriteComments BotPermission = "write:comments" // 写入评论信息
	PermSendMessage      BotPermission = "send:message"   // 发送消息
	PermManageContent    BotPermission = "manage:content" // 管理内容
	PermReadStats        BotPermission = "read:stats"     // 读取统计信息
)

type BotPricing struct {
	Model       BotPricingModel `json:"model" gorm:"type:varchar(20)"`                  // 定价模型
	Price       *float64        `json:"price,omitempty" gorm:"type:decimal(10,2)"`      // 价格
	Cycle       string          `json:"cycle,omitempty" gorm:"type:varchar(20)"`        // 周期
	FreeLimit   string          `json:"freeLimit,omitempty" gorm:"type:text"`           // 免费限制
	PurchaseURL string          `json:"purchaseUrl,omitempty" gorm:"type:varchar(255)"` // 购买链接
}

type BotConfigField struct {
	Key          string             `json:"key"`                    // 唯一标识
	Label        string             `json:"label"`                  // 显示名称
	Type         BotConfigFieldType `json:"type"`                   // 类型
	DefaultValue interface{}        `json:"defaultValue,omitempty"` // 默认值
	Placeholder  string             `json:"placeholder,omitempty"`  // 占位符
	Description  string             `json:"description,omitempty"`  // 描述
	Required     bool               `json:"required"`               // 是否必填
	Options      []struct {
		Label string      `json:"label"` // 显示名称
		Value interface{} `json:"value"` // 值
	} `json:"options,omitempty"` // 选项
}

type BotConfigFieldType string

const (
	FieldTypeText     BotConfigFieldType = "text"     // 文本
	FieldTypeNumber   BotConfigFieldType = "number"   // 数字
	FieldTypeBool     BotConfigFieldType = "bool"     // 布尔
	FieldTypeSelect   BotConfigFieldType = "select"   // 选择
	FieldTypeTextArea BotConfigFieldType = "textarea" // 多行文本
	FieldTypeSecret   BotConfigFieldType = "secret"   // 密码
	FieldTypeArray    BotConfigFieldType = "array"    // 数组
)

type ResourceLimit struct {
	MaxMemoryMB int `json:"maxMemoryMB"` // 最大内存（MB）
	MaxCPU      int `json:"maxCPU"`      // 最大 CPU 核心数
}

func (Bot) TableName() string {
	return "bots"
}

// -----
// 系统默认机器人 ID 常量
const SystemBotID = 1

var DefaultSystemBot = &Bot{
	Name:        "系统助手",
	Version:     "0.1.0",
	Description: "论坛内置机器人，提供基础自动化服务（定时统计、欢迎新用户等）。",
	Summary:     "系统机器人",
	Type:        BotTypeTask,
	CreatorID:   1,
	CreatorName: "System",
	TriggerType: TriggerSchedule,
	CronExpr:    "0 */6 * * *",
	ScriptCode: `-- 系统默认脚本：统计论坛数据
function main()
    local stats = forum.getStats()
    log("当前帖子数: " .. stats.post_count .. ", 用户数: " .. stats.user_count)
end`,
	Enabled:     true,
	Status:      BotStatusActive,
	TimeoutSec:  10,
	Permissions: []BotPermission{PermReadStats},

	// 显式初始化为非 nil JSON 值
	Screenshots: []string{},          // 空 JSON 数组 → []
	Tags:        []string{},          // 空 JSON 数组 → []
	EnvVars:     map[string]string{}, // 空 JSON 对象 → {}
	Pricing: BotPricing{
		Model: BotPricingFree, // 有效枚举值
	},
	ConfigSchema: []BotConfigField{}, // 空 JSON 数组 → []
	ConfigValues: map[string]any{},   // 空 JSON 对象 → {}
}
