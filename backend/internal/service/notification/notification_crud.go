package notification

import (
	"tiny-forum/internal/model"
)

// Create 创建通知（忽略错误，不返回 error 保持原有行为）
func (s *NotificationService) Create(userID uint, senderID *uint, notifType model.NotificationType, content string, targetID *uint, targetType string) {
	n := &model.Notification{
		UserID:     userID,
		SenderID:   senderID,
		Type:       notifType,
		Content:    content,
		TargetID:   targetID,
		TargetType: targetType,
	}
	_ = s.notifRepo.Create(n)
}

// List 获取用户通知列表（分页）
func (s *NotificationService) List(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	return s.notifRepo.ListByUser(userID, page, pageSize)
}

// MarkAllRead 标记所有通知为已读
func (s *NotificationService) MarkAllRead(userID uint) error {
	return s.notifRepo.MarkAllRead(userID)
}

// UnreadCount 获取未读通知数量
func (s *NotificationService) UnreadCount(userID uint) (int64, error) {
	return s.notifRepo.UnreadCount(userID)
}
