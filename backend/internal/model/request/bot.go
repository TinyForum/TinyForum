package request

import "tiny-forum/internal/model/do"

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
