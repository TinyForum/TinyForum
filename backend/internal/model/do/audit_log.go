package do

import "tiny-forum/internal/model/common"

type AuditActionType string

const (
	AuditActionBlockUser      AuditActionType = "block_user"
	AuditActionUnblockUser    AuditActionType = "unblock_user"
	AuditActionHidePost       AuditActionType = "hide_post"
	AuditActionHideComment    AuditActionType = "hide_comment"
	AuditActionHandleReport   AuditActionType = "handle_report"
	AuditActionApproveContent AuditActionType = "approve_content"
	AuditActionRejectContent  AuditActionType = "reject_content"
	AuditActionDeductScore    AuditActionType = "deduct_score"
	AuditActionSetRiskLevel   AuditActionType = "set_risk_level"
)

// AuditLog 管理员操作审计日志，不可删除，只追加
type AuditLog struct {
	common.BaseModel

	OperatorID uint   `gorm:"not null;index" json:"operator_id"`         // 操作者ID
	OperatorIP string `gorm:"type:varchar(45);index" json:"operator_ip"` // 操作者IP

	Action     AuditActionType `gorm:"type:varchar(50);not null;index" json:"action"`                            // 操作类型
	TargetType string          `gorm:"type:varchar(50);not null;index:idx_target,priority:1" json:"target_type"` // 目标类型
	TargetID   uint            `gorm:"not null;index:idx_target,priority:2" json:"target_id"`                    // 目标ID

	Before string `gorm:"type:json" json:"before"` // 操作前数据
	After  string `gorm:"type:json" json:"after"`  // 操作后数据
	Reason string `gorm:"type:text" json:"reason"` // 操作理由
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
