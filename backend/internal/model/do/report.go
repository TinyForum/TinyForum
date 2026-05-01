package do

import (
	"time"
)

// ReportStatus 举报状态枚举
type ReportStatus string

const (
	ReportPending  ReportStatus = "pending"  // 待处理
	ReportResolved ReportStatus = "resolved" // 已处理
	ReportRejected ReportStatus = "rejected" // 已驳回
)

// ReportType 举报类型枚举
type ReportType string

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

// Report 举报记录表
type Report struct {
	BaseModel
	ReporterID uint         `gorm:"column:reporter_id;not null;index;comment:举报人ID" json:"reporter_id"`                              // 举报人ID
	TargetID   uint         `gorm:"column:target_id;not null;index;comment:被举报对象ID" json:"target_id"`                                // 被举报对象ID
	TargetType string       `gorm:"column:target_type;size:50;not null;index;comment:被举报对象类型(post/comment/user)" json:"target_type"` // 被举报对象类型
	Type       ReportType   `gorm:"column:type;type:varchar(50);default:'other';comment:举报类型" json:"type"`                           // 举报类型
	Reason     string       `gorm:"column:reason;size:500;not null;comment:举报理由" json:"reason"`
	Status     ReportStatus `gorm:"column:status;type:varchar(20);default:'pending';index;comment:举报状态" json:"status"` // 举报状态
	HandlerID  *uint        `gorm:"column:handler_id;index;comment:处理人ID" json:"handler_id"`                           // 处理人ID
	HandleNote string       `gorm:"column:handle_note;size:500;comment:处理备注" json:"handle_note"`                       // 处理备注
	HandleAt   *time.Time   `gorm:"column:handle_at;comment:处理时间" json:"handle_at"`                                    // 处理时间

	// 以下为扩展推荐字段（可根据需要启用）
	ContentSnapshot string `gorm:"column:content_snapshot;type:text;comment:被举报内容快照" json:"content_snapshot"` // 被举报内容快照
	ReporterIP      string `gorm:"column:reporter_ip;size:45;comment:举报人IP" json:"reporter_ip"`               // 举报人IP
	IsAnonymous     bool   `gorm:"column:is_anonymous;default:false;comment:是否匿名举报" json:"is_anonymous"`      // 是否匿名举报
	Priority        int8   `gorm:"column:priority;default:2;comment:优先级：1高2中3低" json:"priority"`              // 优先级：1高2中3低

	// 关联对象（查询时通过 Preload 加载，需定义 User 模型）
	Reporter User  `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"` // 举报人
	Handler  *User `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`   // 处理人
}

// TableName 指定表名
func (Report) TableName() string {
	return "reports"
}
