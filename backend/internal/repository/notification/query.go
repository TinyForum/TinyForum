package notification

import "tiny-forum/internal/model/do"

func (r *notificationRepository) ListByUser(userID uint, page, pageSize int) ([]do.Notification, int64, error) {
	var notifications []do.Notification
	var total int64

	query := r.db.Model(&do.Notification{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Sender").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notifications).Error

	return notifications, total, err
}

func (r *notificationRepository) UnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&do.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count).Error
	return count, err
}

func (r *notificationRepository) GetByID(id uint) (*do.Notification, error) {
	var notif do.Notification
	err := r.db.First(&notif, id).Error
	if err != nil {
		return nil, err
	}
	return &notif, nil
}
