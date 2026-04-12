package model

type NotificationType string

const (
	NotifyComment NotificationType = "comment"
	NotifyLike    NotificationType = "like"
	NotifyFollow  NotificationType = "follow"
	NotifyReply   NotificationType = "reply"
	NotifySystem  NotificationType = "system"
)

type Notification struct {
	BaseModel
	UserID     uint             `gorm:"not null;index" json:"user_id"`
	SenderID   *uint            `gorm:"index" json:"sender_id"`
	Type       NotificationType `gorm:"type:varchar(30)" json:"type"`
	Content    string           `gorm:"size:500" json:"content"`
	TargetID   *uint            `json:"target_id"`
	TargetType string           `gorm:"size:50" json:"target_type"`
	IsRead     bool             `gorm:"default:false" json:"is_read"`

	User   User  `gorm:"foreignKey:UserID" json:"-"`
	Sender *User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}
