package notification

import "tiny-forum/internal/model/do"

// BatchMarkRead 批量标记指定通知为已读
func (r *notificationRepository) BatchMarkRead(userID uint, ids []uint) (int64, error) {
	result := r.db.Model(&do.Notification{}).
		Where("user_id = ? AND id IN (?) AND is_read = ?", userID, ids, false).
		Update("is_read", true)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// MarkAllRead 标记所有通知为已读（优化版）
func (r *notificationRepository) MarkAllRead(userID uint) (int64, error) {
	result := r.db.Model(&do.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true)

	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
func (r *notificationRepository) MarkRead(notificationID uint) error {
	return r.db.Model(&do.Notification{}).
		Where("id = ?", notificationID).
		UpdateColumn("is_read", true).Error
}
