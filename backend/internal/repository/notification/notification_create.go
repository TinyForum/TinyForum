package notification

import "tiny-forum/internal/model/do"

func (r *notificationRepository) Create(n *do.Notification) error {
	return r.db.Create(n).Error
}
