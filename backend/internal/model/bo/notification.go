package bo

import "time"

type NotificationBO struct {
	ID         uint
	Type       string
	Content    string
	IsRead     bool
	CreatedAt  time.Time
	TargetID   *uint
	TargetType string
	Sender     *UserBO
}

type UserBO struct {
	ID       uint
	Username string
	Avatar   string
}

// NotificationListResult 分页结果
type NotificationListResult struct {
	List        []*NotificationBO
	Total       int64
	UnreadCount int64
	Page        int
	PageSize    int
}
