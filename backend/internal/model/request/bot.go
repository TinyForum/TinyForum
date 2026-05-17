package request

import (
	"tiny-forum/internal/model/do"
)

type CreateBotRequest struct {
	Name          string              `json:"name" binding:"required,min=1,max=100"`
	Version       string              `json:"version" binding:"required"`
	Description   string              `json:"description,omitempty"`
	Summary       string              `json:"summary,omitempty"`
	AvatarURL     string              `json:"avatarUrl,omitempty" binding:"omitempty,url"`
	Screenshots   []string            `json:"screenshots,omitempty"`
	HomepageURL   string              `json:"homepageUrl,omitempty" binding:"omitempty,url"`
	Type          do.BotType          `json:"type" binding:"required,oneof=chat moderate notify sync task webhook analysis"`
	Tags          []string            `json:"tags,omitempty"`
	ScriptCode    string              `json:"scriptCode" binding:"required_without=ScriptUrl"`
	ScriptURL     string              `json:"scriptUrl,omitempty" binding:"omitempty,url"`
	TriggerType   do.BotTriggerType   `json:"triggerType" binding:"required,oneof=schedule event webhook manual"`
	CronExpr      string              `json:"cronExpr,omitempty"`
	EventFilter   string              `json:"eventFilter,omitempty"`
	TimeoutSec    int                 `json:"timeoutSec" binding:"min=1,max=300"`
	RetryTimes    int                 `json:"retryTimes"`
	EnvVars       map[string]string   `json:"envVars,omitempty"`
	ResourceLimit *do.ResourceLimit   `json:"resourceLimit,omitempty"`
	Pricing       do.BotPricing       `json:"pricing,omitempty"`
	Permissions   []do.BotPermission  `json:"permissions,omitempty"`
	ConfigSchema  []do.BotConfigField `json:"configSchema,omitempty"`
	ConfigValues  map[string]any      `json:"configValues,omitempty"`
}

type UpdateBotRequest struct {
	Name          *string             `json:"name,omitempty"`
	Description   *string             `json:"description,omitempty"`
	Summary       *string             `json:"summary,omitempty"`
	AvatarURL     *string             `json:"avatarUrl,omitempty"`
	ScriptCode    *string             `json:"scriptCode,omitempty"`
	ScriptURL     *string             `json:"scriptUrl,omitempty"`
	TriggerType   *do.BotTriggerType  `json:"triggerType,omitempty"`
	CronExpr      *string             `json:"cronExpr,omitempty"`
	EventFilter   *string             `json:"eventFilter,omitempty"`
	TimeoutSec    *int                `json:"timeoutSec,omitempty"`
	RetryTimes    *int                `json:"retryTimes,omitempty"`
	EnvVars       map[string]string   `json:"envVars,omitempty"`
	ResourceLimit *do.ResourceLimit   `json:"resourceLimit,omitempty"`
	Pricing       *do.BotPricing      `json:"pricing,omitempty"`
	Permissions   []do.BotPermission  `json:"permissions,omitempty"`
	ConfigSchema  []do.BotConfigField `json:"configSchema,omitempty"`
	ConfigValues  map[string]any      `json:"configValues,omitempty"`
	Enabled       *bool               `json:"enabled,omitempty"`
}

// type ValidateFlowRequest struct {
// 	Flow nocode.Flow `json:"flow"`
// }

// ValidateFlowRequest 零代码流程校验请求
type ValidateFlowRequest struct {
	Flow Flow `json:"flow" binding:"required"`
}

// Flow 表示一个完整的零代码机器人流程
type Flow struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Node 表示流程中的一个节点（触发器/条件/动作）
type Node struct {
	ID     string                 `json:"id"`     // 节点唯一标识
	Type   string                 `json:"type"`   // 节点类型，如 "http_trigger", "condition_if", "send_message" 等
	Config map[string]interface{} `json:"config"` // 节点配置，键值对形式
}

// Edge 表示节点之间的连接关系
type Edge struct {
	Source string `json:"source"` // 起始节点 ID
	Target string `json:"target"` // 目标节点 ID
}
