package botapi

import (
	"context"
	"tiny-forum/internal/model/do"
)

// ─── Message ──────────────────────────────────────────────────────────────

// SendMessage 通过 NotificationRepository 创建私信通知。
func (a *forumAPIImpl) SendMessage(ctx context.Context, toUserID uint, content string) error {
	senderID := a.botActorID
	notif := &do.Notification{
		UserID:     uint(toUserID),
		SenderID:   &senderID,
		Type:       do.NotifyPrivateMessage,
		Content:    content,
		TargetType: "bot_message",
	}
	return a.notifRepo.Create(notif)
}
