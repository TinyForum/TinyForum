package do

import (
	"fmt"
	"time"
	"tiny-forum/internal/model/common"
)

// ========================
// 内容审核任务
// ========================

// 审核状态
type ModerationStatus string

const (
	ModerationStatusPending  ModerationStatus = "pending"  // 待审核
	ModerationStatusApproved ModerationStatus = "approved" // 审核通过
	ModerationStatusRejected ModerationStatus = "rejected" // 审核拒绝
)

// 审核类型
type AuditTargetType string

const (
	AuditTargetPost    AuditTargetType = "post"    // 帖子审核
	AuditTargetComment AuditTargetType = "comment" // 评论审核
	AuditTargetUser    AuditTargetType = "user"    // 用户资料审核
)

// 触发类型
type AuditTriggerType string

const (
	AuditTriggerSensitiveWord   AuditTriggerType = "sensitive_word"   // 敏感词触发
	AuditTriggerReportAggregate AuditTriggerType = "report_aggregate" // 举报聚合触发
	AuditTriggerManual          AuditTriggerType = "manual"           // 手动触发
)

// ContentAuditTask 内容审核任务队列
// 命中 review 级敏感词或举报聚合触发后写入，由后台异步处理
type ContentAuditTask struct {
	common.BaseModel
	TargetType  AuditTargetType  `gorm:"type:varchar(20);not null;index" json:"target_type"`     // 审核类型
	TargetID    uint             `gorm:"not null;index" json:"target_id"`                        // 审核目标ID
	TriggerType AuditTriggerType `gorm:"type:varchar(50);not null" json:"trigger_type"`          // 触发类型
	TriggerMeta string           `gorm:"type:text" json:"trigger_meta"`                          // 触发元数据
	Status      ModerationStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"` // 审核状态
	ReviewerID  *uint            `gorm:"index" json:"reviewer_id"`                               // 审核人ID
	ReviewNote  string           `gorm:"size:500" json:"review_note"`                            // 审核备注
	ReviewedAt  *time.Time       `json:"reviewed_at"`                                            // 审核时间

	Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"` // 审核人
}

// 审核任务有效性检验
var validAuditTriggerTypes = map[AuditTriggerType]bool{
	AuditTriggerSensitiveWord:   true,
	AuditTriggerReportAggregate: true,
	AuditTriggerManual:          true,
}

func (at AuditTriggerType) IsValid() bool {
	return validAuditTriggerTypes[at]
}

func (at AuditTriggerType) String() string {
	return string(at)
}

func ParseAuditTriggerType(s string) (AuditTriggerType, error) {
	at := AuditTriggerType(s)
	if !at.IsValid() {
		return "", fmt.Errorf("invalid audit trigger type: %s", s)
	}
	return at, nil
}

// 审核状态有效性验证
var validModerationStatuses = map[ModerationStatus]bool{
	ModerationStatusPending:  true,
	ModerationStatusApproved: true,
	ModerationStatusRejected: true,
}

func (ms ModerationStatus) IsValid() bool {
	return validModerationStatuses[ms]
}
func ParseModerationStatus(s string) ModerationStatus {
	ms := ModerationStatus(s)
	if ms.IsValid() {
		return ms
	}
	return ModerationStatusPending
}
