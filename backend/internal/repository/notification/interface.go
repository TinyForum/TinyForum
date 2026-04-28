package notification

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	// create
	Create(n *model.Notification) error
	// query
	ListByUser(userID uint, page, pageSize int) ([]model.Notification, int64, error)
	GetByID(id uint) (*model.Notification, error)
	UnreadCount(userID uint) (int64, error)
	MarkRead(notiID uint) error
	// update
	MarkAllRead(userID uint) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}
