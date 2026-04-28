package notification

import "tiny-forum/internal/model"

func (r *notificationRepository) MarkAllRead(userID uint) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).
		UpdateColumn("is_read", true).Error
}

func (r *notificationRepository) MarkRead(notificationID uint) error {
    return r.db.Model(&model.Notification{}).
        Where("id = ?", notificationID).
        UpdateColumn("is_read", true).Error
}