package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

type BotStatus string

const (
	BotStatusActive   BotStatus = "active"
	BotStatusInactive BotStatus = "inactive"
	BotStatusError    BotStatus = "error"
	BotStatusLoading  BotStatus = "loading"
	BotStatusStopped  BotStatus = "stopped"
)

type BotType string

const (
	BotTypeChat     BotType = "chat"
	BotTypeModerate BotType = "moderate"
	BotTypeNotify   BotType = "notify"
	BotTypeSync     BotType = "sync"
	BotTypeTask     BotType = "task"
	BotTypeWebhook  BotType = "webhook"
	BotTypeAnalysis BotType = "analysis"
)

type BotTriggerType string

const (
	TriggerSchedule BotTriggerType = "schedule"
	TriggerEvent    BotTriggerType = "event"
	TriggerWebhook  BotTriggerType = "webhook"
	TriggerManual   BotTriggerType = "manual"
)

type BotPricingModel string

const (
	BotPricingFree      BotPricingModel = "free"
	BotPricingFreemium  BotPricingModel = "freemium"
	BotPricingPaid      BotPricingModel = "paid"
	PricingSubscription BotPricingModel = "subscription"
)

type BotPermission string

const (
	BotPermReadUser      BotPermission = "read:user"
	BotPermReadPosts     BotPermission = "read:posts"
	BotPermWritePosts    BotPermission = "write:posts"
	BotPermReadComments  BotPermission = "read:comments"
	BotPermWriteComments BotPermission = "write:comments"
	PermSendMessage      BotPermission = "send:message"
	PermManageContent    BotPermission = "manage:content"
	PermReadStats        BotPermission = "read:stats"
)

type BotPricing struct {
	Model       BotPricingModel `json:"model" gorm:"type:varchar(20)"`
	Price       *float64        `json:"price,omitempty" gorm:"type:decimal(10,2)"`
	Cycle       string          `json:"cycle,omitempty" gorm:"type:varchar(20)"`
	FreeLimit   string          `json:"freeLimit,omitempty" gorm:"type:text"`
	PurchaseURL string          `json:"purchaseUrl,omitempty" gorm:"type:varchar(255)"`
}

type BotConfigField struct {
	Key          string      `json:"key"`
	Label        string      `json:"label"`
	Type         string      `json:"type"` // text, number, boolean, select, textarea, secret
	DefaultValue interface{} `json:"defaultValue,omitempty"`
	Placeholder  string      `json:"placeholder,omitempty"`
	Description  string      `json:"description,omitempty"`
	Required     bool        `json:"required"`
	Options      []struct {
		Label string      `json:"label"`
		Value interface{} `json:"value"`
	} `json:"options,omitempty"`
}

type ResourceLimit struct {
	MaxMemoryMB int `json:"maxMemoryMB"`
	MaxCPU      int `json:"maxCPU"`
}

type Bot struct {
	common.BaseModel
	Name        string   `json:"name" gorm:"type:varchar(100);not null;index"`
	Version     string   `json:"version" gorm:"type:varchar(50);not null"`
	Description string   `json:"description" gorm:"type:text"`
	Summary     string   `json:"summary" gorm:"type:varchar(300)"`
	AvatarURL   string   `json:"avatar_url" gorm:"type:varchar(255)"`
	Screenshots []string `json:"screenshots" gorm:"type:json;serializer:json"` // 添加 serializer
	HomepageURL string   `json:"homepage_url" gorm:"type:varchar(255)"`

	Type BotType  `json:"type" gorm:"type:varchar(30);not null;index"`
	Tags []string `json:"tags" gorm:"type:json;serializer:json"` // 添加 serializer

	CreatorID   uint   `json:"creator_id" gorm:"type:bigint;not null;index"`
	CreatorName string `json:"creator_name" gorm:"type:varchar(100)"`

	ScriptCode string `json:"script_code" gorm:"type:text"`
	ScriptURL  string `json:"script_url" gorm:"type:varchar(500)"`

	TriggerType BotTriggerType `json:"trigger_type" gorm:"type:varchar(20);not null"`
	CronExpr    string         `json:"cron_expr" gorm:"type:varchar(100)"`
	EventFilter string         `json:"event_filter" gorm:"type:varchar(200)"`

	TimeoutSec    int               `json:"timeout_sec" gorm:"default:10"` // 注意 JSON 标签建议用下划线
	RetryTimes    int               `json:"retry_times" gorm:"default:0"`
	EnvVars       map[string]string `json:"env_vars" gorm:"type:json;serializer:json"`       // 添加 serializer
	ResourceLimit *ResourceLimit    `json:"resource_limit" gorm:"type:json;serializer:json"` // 添加 serializer

	Pricing     BotPricing      `json:"pricing" gorm:"type:json;serializer:json"`
	Permissions []BotPermission `json:"permissions" gorm:"type:json;serializer:json"` // 添加 serializer

	Enabled    bool       `json:"enabled" gorm:"default:false;index"`
	Status     BotStatus  `json:"status" gorm:"type:varchar(20);default:'inactive'"`
	ExecCount  int64      `json:"exec_count" gorm:"default:0"`
	LastExecAt *time.Time `json:"last_exec_at" gorm:"type:timestamp"`
	ErrorMsg   string     `json:"error_msg" gorm:"type:text"`

	ConfigSchema []BotConfigField `json:"config_schema" gorm:"type:json;serializer:json"` // 添加 serializer
	ConfigValues map[string]any   `json:"config_values" gorm:"type:json;serializer:json"` // 添加 serializer
}

func (Bot) TableName() string {
	return "bots"
}

// -----
// 系统默认机器人 ID 常量
const SystemBotID = 1

var DefaultSystemBot = &Bot{
	Name:        "系统助手",
	Version:     "1.0.0",
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
