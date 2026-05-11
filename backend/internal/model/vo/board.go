package vo

import "time"

// BoardVO 版块脱敏视图（对外暴露）
type BoardVO struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	Cover       string    `json:"cover,omitempty"`
	ParentID    *uint     `json:"parent_id,omitempty"` // 仅保留父版块ID，不嵌套完整对象
	SortOrder   int       `json:"sort_order"`
	ViewRole    string    `json:"view_role"` // UserRole 映射为字符串
	PostRole    string    `json:"post_role"`
	ReplyRole   string    `json:"reply_role"`
	PostCount   int       `json:"post_count"`
	ThreadCount int       `json:"thread_count"`
	TodayCount  int       `json:"today_count"`
}
