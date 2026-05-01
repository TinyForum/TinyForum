package notification

import (
	"tiny-forum/internal/model/po"
)

func (r *notificationRepository) Create(n *po.Notification) error {
	return r.db.Create(n).Error
}
