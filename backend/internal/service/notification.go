// package service

// import (
// 	"tiny-forum/internal/model"
// 	"tiny-forum/internal/repository"
// )

// type NotificationService struct {
// 	notifRepo *repository.NotificationRepository
// }

// func NewNotificationService(notifRepo *repository.NotificationRepository) *NotificationService {
// 	return &NotificationService{notifRepo: notifRepo}
// }

// func (s *NotificationService) Create(userID uint, senderID *uint, notifType model.NotificationType, content string, targetID *uint, targetType string) {
// 	n := &model.Notification{
// 		UserID:     userID,
// 		SenderID:   senderID,
// 		Type:       notifType,
// 		Content:    content,
// 		TargetID:   targetID,
// 		TargetType: targetType,
// 	}
// 	_ = s.notifRepo.Create(n)
// }

// func (s *NotificationService) List(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
// 	return s.notifRepo.ListByUser(userID, page, pageSize)
// }

// func (s *NotificationService) MarkAllRead(userID uint) error {
// 	return s.notifRepo.MarkAllRead(userID)
// }

// func (s *NotificationService) UnreadCount(userID uint) (int64, error) {
// 	return s.notifRepo.UnreadCount(userID)
// }

package service
