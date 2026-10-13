package vo

import "time"

// AuditTaskVO 审核任务视图对象
type AuditTaskVO struct {
	ID          uint      `json:"id"`           // 任务ID
	CreatedAt   time.Time `json:"created_at"`   // 创建时间
	TargetType  string    `json:"target_type"`  // 审核目标类型：post / comment / user
	TargetID    uint      `json:"target_id"`    // 审核目标ID
	TriggerType string    `json:"trigger_type"` // 触发类型
	TriggerMeta string    `json:"trigger_meta"` // 触发元数据
	Status      string    `json:"status"`       // 审核状态

	// 以下字段仅在审核历史中返回
	ReviewerID   *uint      `json:"reviewer_id,omitempty"`   // 审核人ID
	ReviewerName string     `json:"reviewer_name,omitempty"` // 审核人用户名
	ReviewNote   string     `json:"review_note,omitempty"`   // 审核备注
	ReviewedAt   *time.Time `json:"reviewed_at,omitempty"`   // 审核时间
}
