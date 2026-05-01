package notification

import "tiny-forum/internal/model/po"

func (r *notificationRepository) ListByUser(userID uint, page, pageSize int) ([]po.Notification, int64, error) {
	var notifications []po.Notification
	var total int64

	query := r.db.Model(&po.Notification{}).Where("user_id = ?", userID)
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
	err := r.db.Model(&po.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count).Error
	return count, err
}

func (r *notificationRepository) GetByID(id uint) (*po.Notification, error) {
	var notif po.Notification
	err := r.db.First(&notif, id).Error
	if err != nil {
		return nil, err
	}
	return &notif, nil
}
