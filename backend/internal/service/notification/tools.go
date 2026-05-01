// service/notification/converter.go

package notification

import (
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/po"
)

// modelToBO 单个 Model 转 BO
func modelToBO(notif *po.Notification) *bo.NotificationBO {
	if notif == nil {
		return nil
	}

	notificationBO := &bo.NotificationBO{
		ID:         notif.ID,
		Type:       string(notif.Type),
		Content:    notif.Content,
		IsRead:     notif.IsRead,
		CreatedAt:  notif.CreatedAt,
		TargetID:   notif.TargetID,
		TargetType: notif.TargetType,
	}

	if notif.Sender != nil {
		notificationBO.Sender = &bo.UserBO{
			ID:       notif.Sender.ID,
			Username: notif.Sender.Username,
			Avatar:   notif.Sender.Avatar,
		}
	}

	return notificationBO
}

// modelsToBOs 批量 Model 转 BO
func modelsToBOs(notifications []po.Notification) []*bo.NotificationBO {
	bos := make([]*bo.NotificationBO, 0, len(notifications))
	for i := range notifications {
		bos = append(bos, modelToBO(&notifications[i]))
	}
	return bos
}
