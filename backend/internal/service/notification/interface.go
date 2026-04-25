package notification

import (
	"tiny-forum/internal/model"
	notificationRepo "tiny-forum/internal/repository/notification"
)

type NotificationService interface {
	Create(userID uint, senderID *uint, notifType model.NotificationType, content string, targetID *uint, targetType string)
	List(userID uint, page, pageSize int) ([]model.Notification, int64, error)
	MarkAllRead(userID uint) error
	UnreadCount(userID uint) (int64, error)
}
type notificationService struct {
	notifRepo notificationRepo.NotificationRepository
}

func NewNotificationService(notifRepo notificationRepo.NotificationRepository) NotificationService {
	return &notificationService{notifRepo: notifRepo}
}
