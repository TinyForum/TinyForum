package po

import (
	"database/sql"
	"time"
)

// ViolationType 违规类型枚举
type ViolationType int8

const (
	ViolationTypeAbuse     ViolationType = 1 // 辱骂
	ViolationTypeAd        ViolationType = 2 // 广告
	ViolationTypeSpam      ViolationType = 3 // 刷屏
	ViolationTypeFraud     ViolationType = 4 // 欺诈
	ViolationTypeSensitive ViolationType = 5 // 敏感内容
)

// SourceType 来源枚举
type SourceType int8

const (
	SourceSystem   SourceType = 1 // 系统自动
	SourceAdmin    SourceType = 2 // 管理员
	SourceUserFlag SourceType = 3 // 用户举报
)

// ViolationStatus 状态枚举
type ViolationStatus int8

const (
	ViolationStatusPending   ViolationStatus = 1 // 待处理
	ViolationStatusConfirmed ViolationStatus = 2 // 已确认
	ViolationStatusRejected  ViolationStatus = 3 // 驳回
)

// PunishType 处罚类型枚举
type PunishType int8

const (
	PunishTypeWarning  PunishType = 1 // 警告
	PunishTypeMute     PunishType = 2 // 禁言
	PunishTypeBan      PunishType = 3 // 封禁
	PunishTypeRestrict PunishType = 4 // 限制功能
)

// AppealStatus 申诉状态枚举
type AppealStatus int8

const (
	AppealStatusNone     AppealStatus = 0 // 未申诉
	AppealStatusPending  AppealStatus = 1 // 申诉中
	AppealStatusApproved AppealStatus = 2 // 申诉通过
	AppealStatusRejected AppealStatus = 3 // 申诉驳回
)

type Violation struct {
	BaseModel
	UserID          int64           `gorm:"column:user_id;not null;index:idx_user_id;comment:违规用户ID" json:"user_id"`
	ViolationType   ViolationType   `gorm:"column:violation_type;not null;comment:违规类型" json:"violation_type"`
	Reason          string          `gorm:"column:reason;size:255;not null;default:'';comment:违规原因简述" json:"reason"`
	ContentSnapshot sql.NullString  `gorm:"column:content_snapshot;type:text;comment:违规内容快照(JSON或文本)" json:"content_snapshot"`
	EvidenceURL     string          `gorm:"column:evidence_url;size:512;default:'';comment:截图或证据文件URL" json:"evidence_url"`
	Source          SourceType      `gorm:"column:source;not null;default:1;comment:来源" json:"source"`
	Status          ViolationStatus `gorm:"column:status;not null;default:1;comment:处理状态" json:"status"`
	OperatorID      int64           `gorm:"column:operator_id;default:0;comment:操作人ID(管理员)" json:"operator_id"`
	PunishType      PunishType      `gorm:"column:punish_type;default:0;comment:处罚类型" json:"punish_type"`
	PunishExpireAt  *time.Time      `gorm:"column:punish_expire_at;comment:处罚过期时间(临时处罚)" json:"punish_expire_at"`
	AppealStatus    AppealStatus    `gorm:"column:appeal_status;not null;default:0;comment:申诉状态" json:"appeal_status"`
	AppealReason    string          `gorm:"column:appeal_reason;size:500;default:'';comment:申诉理由" json:"appeal_reason"`
	AppealTime      *time.Time      `gorm:"column:appeal_time;comment:申诉时间" json:"appeal_time"`
	AppealResult    string          `gorm:"column:appeal_result;size:500;default:'';comment:申诉处理结果" json:"appeal_result"`
}

// TableName 指定表名
func (Violation) TableName() string {
	return "violations"
}
