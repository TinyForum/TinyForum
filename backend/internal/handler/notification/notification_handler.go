package notification

import (
	notificationService "tiny-forum/internal/service/notification"
)

type NotificationHandler struct {
	notifSvc notificationService.NotificationService
}

func NewNotificationHandler(notifSvc notificationService.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifSvc: notifSvc}
}
