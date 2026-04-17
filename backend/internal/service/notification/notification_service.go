package notification

import (
	notificationRepo "tiny-forum/internal/repository/notification"
)

type NotificationService struct {
	notifRepo *notificationRepo.NotificationRepository
}

func NewNotificationService(notifRepo *notificationRepo.NotificationRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo}
}
