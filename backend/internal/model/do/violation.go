package do

import (
	"time"
	"tiny-forum/internal/model/common"
	apperrors "tiny-forum/pkg/errors"
)

// ========================
// Violation 违规记录表
// ========================

// Violation 记录用户违规行为及处理结果
type Violation struct {
	common.BaseModel
	UserID          uint            `gorm:"not null;index:idx_user_violation,priority:1;comment:违规用户ID" json:"user_id"`   // 违规用户ID
	ViolationType   ViolationType   `gorm:"type:varchar(30);not null;index;comment:违规类型" json:"violation_type"`           // 违规类型
	Reason          string          `gorm:"type:varchar(255);not null;comment:违规原因简述" json:"reason"`                      // 违规原因简述
	ContentSnapshot string          `gorm:"type:text;comment:违规内容快照(JSON或纯文本)" json:"content_snapshot"`                   // 违规内容快照
	EvidenceURL     string          `gorm:"type:varchar(512);comment:证据文件URL" json:"evidence_url"`                        // 证据文件URL
	Source          VoilationSource `gorm:"type:varchar(20);not null;default:'system';comment:来源" json:"source"`          // 来源
	Status          ViolationStatus `gorm:"type:varchar(20);not null;default:'pending';index;comment:处理状态" json:"status"` // 处理状态
	OperatorID      uint            `gorm:"default:0;comment:操作人ID(管理员)" json:"operator_id"`                              // 操作人ID
	PunishType      PunishType      `gorm:"type:varchar(30);default:'';comment:处罚类型" json:"punish_type"`                  // 处罚类型
	PunishExpireAt  *time.Time      `gorm:"comment:处罚过期时间(临时处罚)" json:"punish_expire_at"`                                 // 处罚过期时间
	AppealStatus    AppealStatus    `gorm:"type:varchar(20);not null;default:'none';comment:申诉状态" json:"appeal_status"`   // 申诉状态
	AppealReason    string          `gorm:"type:varchar(500);comment:申诉理由" json:"appeal_reason"`                          // 申诉理由
	AppealTime      *time.Time      `gorm:"comment:申诉时间" json:"appeal_time"`                                              // 申诉时间
	AppealResult    string          `gorm:"type:varchar(500);comment:申诉处理结果" json:"appeal_result"`                        // 申诉处理结果

	// 关联（可选）
	User     User `gorm:"foreignKey:UserID" json:"user,omitempty"`         // 违规用户
	Operator User `gorm:"foreignKey:OperatorID" json:"operator,omitempty"` // 操作人
}

// TableName 指定表名（GORM 默认会使用复数 violations，此处显式声明可选）
func (Violation) TableName() string {
	return "violations"
}

// ViolationType 违规类型
type ViolationType string

const (
	ViolationTypeAbuse     ViolationType = "abuse"     // 辱骂
	ViolationTypeAd        ViolationType = "ad"        // 广告
	ViolationTypeSpam      ViolationType = "spam"      // 刷屏
	ViolationTypeFraud     ViolationType = "fraud"     // 欺诈
	ViolationTypeSensitive ViolationType = "sensitive" // 敏感内容
)

// enum [abuse, ad, spam, fraud, sensitive]

// IsValid 验证违规类型是否合法
func (vt ViolationType) IsValid() bool {
	switch vt {
	case ViolationTypeAbuse, ViolationTypeAd, ViolationTypeSpam,
		ViolationTypeFraud, ViolationTypeSensitive:
		return true
	}
	return false
}

// 解析违规类型
func ParseViolationType(s string) (ViolationType, error) {
	vt := ViolationType(s)
	if vt.IsValid() {
		return vt, nil
	}
	return "", apperrors.ErrValidation
}

// ViolationStatus 违规处理状态
type ViolationStatus string

const (
	ViolationStatusPending   ViolationStatus = "pending"   // 待处理
	ViolationStatusConfirmed ViolationStatus = "confirmed" // 已确认
	ViolationStatusRejected  ViolationStatus = "rejected"  // 驳回
)

// enum [pending, confirmed, rejected]

// 解析违规处理状态
func ParseViolationStatus(s string) (ViolationStatus, error) {
	vs := ViolationStatus(s)
	if vs.IsValid() {
		return vs, nil
	}
	return "", apperrors.ErrValidation
}

func (vs ViolationStatus) IsValid() bool {
	switch vs {
	case ViolationStatusPending, ViolationStatusConfirmed, ViolationStatusRejected:
		return true
	}
	return false
}

// PunishType 处罚类型
type PunishType string

const (
	PunishTypeWarning  PunishType = "warning"  // 警告
	PunishTypeMute     PunishType = "mute"     // 禁言
	PunishTypeBan      PunishType = "ban"      // 封禁
	PunishTypeRestrict PunishType = "restrict" // 限制功能
)

// enum [warning, mute, ban, restrict]

func ParsePunishType(s string) (PunishType, error) {
	pt := PunishType(s)
	if pt.IsValid() {
		return pt, nil
	}
	return "", apperrors.ErrValidation
}
func (pt PunishType) IsValid() bool {
	switch pt {
	case PunishTypeWarning, PunishTypeMute, PunishTypeBan, PunishTypeRestrict:
		return true
	}
	return false
}

// AppealStatus 申诉状态
type AppealStatus string

const (
	AppealStatusNone     AppealStatus = "none"     // 未申诉
	AppealStatusPending  AppealStatus = "pending"  // 申诉中
	AppealStatusApproved AppealStatus = "approved" // 申诉通过
	AppealStatusRejected AppealStatus = "rejected" // 申诉驳回
)

// enum [none, pending, approved, rejected]

func ParseAppealStatus(s string) (AppealStatus, error) {
	as := AppealStatus(s)
	if as.IsValid() {
		return as, nil
	}
	return "", apperrors.ErrValidation
}

func (as AppealStatus) IsValid() bool {
	switch as {
	case AppealStatusNone, AppealStatusPending, AppealStatusApproved, AppealStatusRejected:
		return true
	}
	return false
}

// 违规来源
type VoilationSource string

const (
	VoilationSourceUser     VoilationSource = "user"
	VoilationSourceAdmin    VoilationSource = "admin"
	VoilationSourceReviewer VoilationSource = "reviewer"
	VoilationSourceSystem   VoilationSource = "system"
)

// enum [user, admin, reviewer, system]

func ParseVoilationSource(s string) (VoilationSource, error) {
	vs := VoilationSource(s)
	if vs.IsValid() {
		return vs, nil
	}
	return "", apperrors.ErrValidation
}
func (st VoilationSource) IsValid() bool {
	switch st {
	case VoilationSourceUser, VoilationSourceAdmin, VoilationSourceReviewer, VoilationSourceSystem:
		return true
	}
	return false
}
