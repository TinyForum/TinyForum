package notification

import (
	"tiny-forum/internal/model"
)

func (r *NotificationRepository) Create(n *model.Notification) error {
	return r.db.Create(n).Error
}
