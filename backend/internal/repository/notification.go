package repository

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(n *model.Notification) error {
	return r.db.Create(n).Error
}

func (r *NotificationRepository) ListByUser(userID uint, page, pageSize int) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64
	query := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Preload("Sender").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications).Error
	return notifications, total, err
}

func (r *NotificationRepository) MarkAllRead(userID uint) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).
		UpdateColumn("is_read", true).Error
}

func (r *NotificationRepository) UnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = false", userID).Count(&count).Error
	return count, err
}
