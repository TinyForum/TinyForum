package model

type ReportStatus string

const (
	ReportPending  ReportStatus = "pending"
	ReportResolved ReportStatus = "resolved"
	ReportRejected ReportStatus = "rejected"
)

type Report struct {
	BaseModel
	ReporterID uint         `gorm:"not null;index" json:"reporter_id"`
	TargetID   uint         `gorm:"not null" json:"target_id"`
	TargetType string       `gorm:"size:50;not null" json:"target_type"` // post | comment | user
	Reason     string       `gorm:"size:500;not null" json:"reason"`
	Status     ReportStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	HandlerID  *uint        `json:"handler_id"`
	HandleNote string       `gorm:"size:500" json:"handle_note"`

	Reporter User  `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	Handler  *User `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`
}
