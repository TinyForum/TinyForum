package po

import "time"

// ========================
// 用户风险等级
// ========================

type RiskLevel string

const (
	RiskLevelNormal   RiskLevel = "normal"   // 正常用户
	RiskLevelObserve  RiskLevel = "observe"  // 观察中（积分低 / 新用户）
	RiskLevelRestrict RiskLevel = "restrict" // 受限（近期有被处理的举报）
	RiskLevelBlocked  RiskLevel = "blocked"  // 封禁（对应 User.IsBlocked）
)

// ========================
// 内容审核任务
// ========================

type ModerationStatus string

const (
	ModerationStatusPending  ModerationStatus = "pending"  // 待审核
	ModerationStatusApproved ModerationStatus = "approved" // 审核通过
	ModerationStatusRejected ModerationStatus = "rejected" // 审核拒绝（内容被隐藏）
)

// 审核类型
type AuditTargetType string

const (
	AuditTargetPost    AuditTargetType = "post"
	AuditTargetComment AuditTargetType = "comment"
	AuditTargetUser    AuditTargetType = "user" // 用户资料审核（bio/username）
)

// ContentAuditTask 内容审核任务队列
// 命中 review 级敏感词或举报聚合触发后写入，由后台异步处理
type ContentAuditTask struct {
	BaseModel
	TargetType  AuditTargetType  `gorm:"type:varchar(20);not null;index" json:"target_type"`
	TargetID    uint             `gorm:"not null;index" json:"target_id"`
	TriggerType string           `gorm:"type:varchar(50);not null" json:"trigger_type"` // "sensitive_word" | "report_aggregate" | "manual"
	TriggerMeta string           `gorm:"type:text" json:"trigger_meta"`                 // JSON：命中的词、举报数等
	Status      ModerationStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	ReviewerID  *uint            `gorm:"index" json:"reviewer_id"`
	ReviewNote  string           `gorm:"size:500" json:"review_note"`
	ReviewedAt  *time.Time       `json:"reviewed_at"`

	Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// ========================
// 操作审计日志
// ========================

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
	BaseModel
	OperatorID uint            `gorm:"not null;index" json:"operator_id"`
	OperatorIP string          `gorm:"type:varchar(45);index"` // 操作者IP
	Action     AuditActionType `gorm:"type:varchar(50);not null;index" json:"action"`
	TargetType string          `gorm:"type:varchar(50);not null" json:"target_type"`
	TargetID   uint            `gorm:"not null" json:"target_id"`
	Before     string          `gorm:"type:text" json:"before"` // JSON 序列化的变更前状态
	After      string          `gorm:"type:text" json:"after"`  // JSON 序列化的变更后状态
	Reason     string          `gorm:"size:500" json:"reason"`
	IP         string          `gorm:"size:64" json:"ip"`

	Operator User `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// ========================
// 风控记录
// ========================

// UserRiskRecord 记录用户风险事件，用于计算当前风险等级
type UserRiskRecord struct {
	BaseModel
	UserID uint `gorm:"not null;index" json:"user_id"` // 用户ID
	// "report_confirmed" | "sensitive_hit" | "rate_limit_exceeded"
	EventType   string    `gorm:"type:varchar(50);not null" json:"event_type"` // 事件类型
	EventDetail string    `gorm:"type:text" json:"event_detail"`               // 事件详情
	ExpireAt    time.Time `gorm:"not null;index" json:"expire_at"`             // 超过此时间后不计入风险分

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 关联用户
}

func (UserRiskRecord) TableName() string {
	return "user_risk_records"
}

// IPRiskRecord IP风险记录
type IPRiskRecord struct {
	ID          uint      `gorm:"primaryKey"`
	IP          string    `gorm:"type:varchar(45);index;not null"`
	EventType   string    `gorm:"type:varchar(50);not null"`
	EventDetail string    `gorm:"type:text"`
	ExpireAt    time.Time `gorm:"index"`
	CreatedAt   time.Time
}

func (IPRiskRecord) TableName() string {
	return "ip_risk_records"
}

// BlockedIP IP封禁记录
type BlockedIP struct {
	ID         uint       `gorm:"primaryKey"`
	IP         string     `gorm:"type:varchar(45);uniqueIndex;not null"` // 支持IPv6
	Reason     string     `gorm:"type:text"`                             // 封禁原因
	OperatorID uint       `gorm:"index"`                                 // 操作员ID
	ExpireAt   *time.Time `gorm:"index"`                                 // 封禁过期时间，NULL表示永久封禁
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (BlockedIP) TableName() string {
	return "blocked_ips"
}

// IsExpired 检查封禁是否已过期
func (b *BlockedIP) IsExpired() bool {
	if b.ExpireAt == nil {
		return false
	}
	return time.Now().After(*b.ExpireAt)
}

// IsPermanent 是否是永久封禁
func (b *BlockedIP) IsPermanent() bool {
	return b.ExpireAt == nil
}
