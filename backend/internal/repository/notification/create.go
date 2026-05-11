package notification

import (
	"tiny-forum/internal/model/do"
)

func (r *notificationRepository) Create(n *do.Notification) error {
	return r.db.Create(n).Error
}

// func (r *notificationRepository) Create(ctx context.Context, userID uint, sendID *uint, notifyType do.NotificationType, targetID *uint, targetType string, content string, IsRead bool) error {
// 	n := &do.Notification{
// 		UserID:     userID,
// 		SenderID:   sendID,
// 		Type:       notifyType,
// 		Content:    content,
// 		TargetID:   targetID,
// 		TargetType: targetType,
// 		IsRead:     IsRead,
// 	}
// 	return r.db.Create(n).Error
// }
