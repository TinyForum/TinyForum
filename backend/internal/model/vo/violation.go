package vo

import "time"

// ViolationVO 违规记录脱敏视图（对外暴露）
type ViolationVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID     int64 `json:"user_id"`
	OperatorID int64 `json:"operator_id,omitempty"`

	ViolationType string `json:"violation_type"` // 违规类型
	Reason        string `json:"reason"`         // 违规原因简述
	Source        string `json:"source"`         // 来源（system/user/admin等）
	Status        string `json:"status"`         // 处理状态

	PunishType     string     `json:"punish_type,omitempty"`      // 处罚类型
	PunishExpireAt *time.Time `json:"punish_expire_at,omitempty"` // 临时处罚过期时间

	AppealStatus string     `json:"appeal_status"`           // 申诉状态
	AppealReason string     `json:"appeal_reason,omitempty"` // 申诉理由
	AppealTime   *time.Time `json:"appeal_time,omitempty"`   // 申诉时间
	AppealResult string     `json:"appeal_result,omitempty"` // 申诉处理结果
}
