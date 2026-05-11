package vo

import "time"

// NotificationVO 通知脱敏视图（对外暴露）
type NotificationVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID     uint   `json:"user_id"` // 接收通知的用户ID
	SenderID   *uint  `json:"sender_id,omitempty"`
	Type       string `json:"type"` // NotificationType 映射为字符串
	Content    string `json:"content"`
	TargetID   *uint  `json:"target_id,omitempty"`
	TargetType string `json:"target_type,omitempty"`
	IsRead     bool   `json:"is_read"`

	// 脱敏后的发送者信息（需额外查询填充）
	Sender struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar,omitempty"`
	} `json:"sender,omitempty"`
}
