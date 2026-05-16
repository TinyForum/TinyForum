package do

import (
	"fmt"
	"time"
	"tiny-forum/internal/model/common"
)

// UserRiskRecord 用户风险记录
type UserRiskRecord struct {
	common.BaseModel
	UserID      uint          `gorm:"not null;index:idx_user_expire,priority:1" json:"user_id"`   // 用户ID
	EventType   RiskEventType `gorm:"type:varchar(50);not null" json:"event_type"`                // 事件类型
	EventDetail string        `gorm:"type:text" json:"event_detail"`                              // 事件详情
	RiskLevel   RiskLevel     `gorm:"type:varchar(50);not null" json:"risk_level"`                // 风险等级s
	ExpireAt    time.Time     `gorm:"not null;index:idx_user_expire,priority:2" json:"expire_at"` // 过期时间

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"` // 关联用户
}

func (UserRiskRecord) TableName() string {
	return "user_risk_records"
}

// RiskLevel 用户风险等级
type RiskLevel string

const (
	RiskLevelNormal   RiskLevel = "normal"   // 正常用户
	RiskLevelObserve  RiskLevel = "observe"  // 观察中
	RiskLevelRestrict RiskLevel = "restrict" // 受限
	RiskLevelBlocked  RiskLevel = "blocked"  // 封禁（对应 User.IsBlocked）
)

// enum [normal,observe,restrict,blocked]

// RiskEventType 风险事件类型
type RiskEventType string

const (
	RiskEventTypeReportConfirmed   RiskEventType = "report_confirmed"    // 举报确认
	RiskEventTypeSensitiveHit      RiskEventType = "sensitive_hit"       // 敏感词命中
	RiskEventTypeRateLimitExceeded RiskEventType = "rate_limit_exceeded" // 请求频率过高
)

// enum [report_confirmed,sensitive_hit,rate_limit_exceeded]

// IsValid 验证事件类型
func (t RiskEventType) IsValid() bool {
	switch t {
	case RiskEventTypeReportConfirmed,
		RiskEventTypeSensitiveHit,
		RiskEventTypeRateLimitExceeded:
		return true
	}
	return false
}

// ParseRiskEventType 严格解析，返回错误
func ParseRiskEventType(s string) (RiskEventType, error) {
	t := RiskEventType(s)
	if t.IsValid() {
		return t, nil
	}
	return "", fmt.Errorf("invalid risk event type: %s", s)
}

// IsExpired 判断记录是否已过期
func (r *UserRiskRecord) IsExpired(now time.Time) bool {
	return !r.ExpireAt.After(now)
}
