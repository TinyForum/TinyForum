package notification

import "tiny-forum/internal/model"

func (r *NotificationRepository) MarkAllRead(userID uint) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).
		UpdateColumn("is_read", true).Error
}
