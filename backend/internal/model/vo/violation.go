package vo

import (
	"time"
	"tiny-forum/internal/model/do"
)

// ViolationVO 违规记录视图对象（用于 API 响应）
type ViolationVO struct {
	ID             uint               `json:"id"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	UserID         uint               `json:"user_id"`
	OperatorID     *uint              `json:"operator_id,omitempty"` // 操作人ID，可为空
	ViolationType  do.ViolationType   `json:"violation_type"`
	Reason         string             `json:"reason"`
	Source         do.VoilationSource `json:"source"`
	Status         do.ViolationStatus `json:"status"`
	PunishType     do.PunishType      `json:"punish_type,omitempty"`
	PunishExpireAt *time.Time         `json:"punish_expire_at,omitempty"`
	AppealStatus   do.AppealStatus    `json:"appeal_status"`
	AppealReason   string             `json:"appeal_reason,omitempty"`
	AppealTime     *time.Time         `json:"appeal_time,omitempty"`
	AppealResult   string             `json:"appeal_result,omitempty"`
	User           *UserAuthVO        `json:"user,omitempty"`
	Operator       *UserAuthVO        `json:"operator,omitempty"`
}
