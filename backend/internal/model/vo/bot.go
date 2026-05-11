package vo

import (
	"time"
	"tiny-forum/internal/model/do"
)

type BotResponse struct {
	ID            uint                `json:"id"`
	Name          string              `json:"name"`
	Version       string              `json:"version"`
	Description   string              `json:"description"`
	Summary       string              `json:"summary"`
	AvatarURL     string              `json:"avatarUrl"`
	Screenshots   []string            `json:"screenshots"`
	HomepageURL   string              `json:"homepageUrl"`
	Type          do.BotType          `json:"type"`
	Tags          []string            `json:"tags"`
	CreatorID     uint                `json:"creatorId"`
	CreatorName   string              `json:"creatorName"`
	TriggerType   do.BotTriggerType   `json:"triggerType"`
	CronExpr      string              `json:"cronExpr"`
	EventFilter   string              `json:"eventFilter"`
	TimeoutSec    int                 `json:"timeoutSec"`
	RetryTimes    int                 `json:"retryTimes"`
	ResourceLimit *do.ResourceLimit   `json:"resourceLimit"`
	Pricing       do.BotPricing       `json:"pricing"`
	Permissions   []do.BotPermission  `json:"permissions"`
	Enabled       bool                `json:"enabled"`
	Status        do.BotStatus        `json:"status"`
	ExecCount     int64               `json:"execCount"`
	LastExecAt    *time.Time          `json:"lastExecAt"`
	ErrorMsg      string              `json:"errorMsg"`
	ConfigSchema  []do.BotConfigField `json:"configSchema,omitempty"`
	ConfigValues  map[string]any      `json:"configValues,omitempty"`
	CreatedAt     time.Time           `json:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt"`
}

type BotExecutionLogResponse struct {
	ID        uint      `json:"id"`
	BotID     uint      `json:"botId"`
	Status    string    `json:"status"`   // success, error
	Duration  int64     `json:"duration"` // ms
	Output    string    `json:"output"`
	ErrorMsg  string    `json:"errorMsg"`
	CreatedAt time.Time `json:"createdAt"`
}
