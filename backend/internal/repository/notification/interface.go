package notification

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	// create
	Create(n *do.Notification) error
	// query
	ListByUser(userID uint, page, pageSize int) ([]do.Notification, int64, error)
	GetByID(id uint) (*do.Notification, error)
	UnreadCount(userID uint) (int64, error)
	MarkRead(notiID uint) error
	// update
	MarkAllRead(userID uint) (int64, error)
	BatchMarkRead(userID uint, ids []uint) (int64, error) // 新增批量标记

}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}
