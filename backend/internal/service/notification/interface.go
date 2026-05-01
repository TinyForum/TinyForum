package notification

import (
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/do"
	notificationRepo "tiny-forum/internal/repository/notification"
)

type NotificationService interface {
	Create(userID uint, senderID *uint, notifType do.NotificationType, content string, targetID *uint, targetType string)
	List(userID uint, page, pageSize int) (*bo.NotificationListResult, error)
	MarkAllRead(userID uint) (int64, error)
	MarkRead(userID uint, notifID uint) error
	UnreadCount(userID uint) (int64, error)
	BatchMarkRead(userID uint, ids []uint) (int64, error) // 新增批量标记
}
type notificationService struct {
	notifRepo notificationRepo.NotificationRepository
}

func NewNotificationService(notifRepo notificationRepo.NotificationRepository) NotificationService {
	return &notificationService{notifRepo: notifRepo}
}
