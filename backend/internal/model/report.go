package model

import "time"

type ReportStatus string
type ReportType string

const (
	ReportPending  ReportStatus = "pending"
	ReportResolved ReportStatus = "resolved"
	ReportRejected ReportStatus = "rejected"
)

const (
	ReportTypeSpam           ReportType = "spam"           // 广告/垃圾信息
	ReportTypeOffensive      ReportType = "offensive"      // 侮辱/攻击性内容
	ReportTypeIllegal        ReportType = "illegal"        // 违法内容
	ReportTypeMisinformation ReportType = "misinformation" // 虚假信息
	ReportTypePrivacy        ReportType = "privacy"        // 侵犯隐私
	ReportTypeOther          ReportType = "other"          // 其他
)

// ReportAggregateThreshold 同一内容被举报多少次后自动进入审核队列
const ReportAggregateThreshold = 3

type Report struct {
	BaseModel
	ReporterID uint         `gorm:"not null;index" json:"reporter_id"`
	TargetID   uint         `gorm:"not null;index" json:"target_id"`
	TargetType string       `gorm:"size:50;not null;index" json:"target_type"` // post | comment | user
	Type       ReportType   `gorm:"type:varchar(50);default:'other'" json:"type"`
	Reason     string       `gorm:"size:500;not null" json:"reason"`
	Status     ReportStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`
	HandlerID  *uint        `json:"handler_id"`
	HandleNote string       `gorm:"size:500" json:"handle_note"`
	HandleAt   *time.Time   `json:"handle_at"` // 处理时间，用于统计处理效率

	Reporter User  `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	Handler  *User `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`
}
