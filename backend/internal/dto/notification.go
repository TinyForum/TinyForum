package dto

import "time"

type NotificationListRequest struct {
	Page     int `form:"page" json:"page" binding:"min=1"`
	PageSize int `form:"page_size" json:"page_size" binding:"min=1,max=100"`
}
type NotificationSenderResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar,omitempty"`
}

type NotificationResponse struct {
	ID        uint                       `json:"id"`
	Type      string                     `json:"type"`
	Content   string                     `json:"content"`
	IsRead    bool                       `json:"is_read"`
	CreatedAt time.Time                  `json:"created_at"`
	Sender    *NotificationSenderResponse `json:"sender,omitempty"`
	TargetID  *uint                      `json:"target_id,omitempty"`
	TargetType string                     `json:"target_type,omitempty"`
}
// 
type NotificationListResponse struct {
	List       []NotificationResponse `json:"list"`
	Total      int64                  `json:"total"`
	UnreadCount int64                 `json:"unread_count"`
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
}



// BatchMarkReadRequest 批量标记已读请求
type BatchMarkReadRequest struct {
	IDs []uint `json:"ids"` // 通知ID列表，为空则标记所有
}

// BatchMarkReadResponse 批量标记已读响应
type BatchMarkReadResponse struct {
	Message      string `json:"message"`
	UpdatedCount int64    `json:"updated_count"` // 实际更新的数量
}