package botapi

import (
	"context"
	"fmt"
	"time"
	"tiny-forum/internal/infra/lua/sdk"
	"tiny-forum/internal/model/do"
)

// ─── User ─────────────────────────────────────────────────────────────────

func (a *forumAPIImpl) GetUser(ctx context.Context, userID uint) (*sdk.UserVO, error) {
	u, err := a.userRepo.FindByID(uint(userID))
	if err != nil {
		return nil, err
	}
	return &sdk.UserVO{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Role:     string(u.Role),
	}, nil
}

// BanUser 通过 UserRepository.UpdateBlocked 封禁用户，并发送通知。
// durationSec 目前实现为永久封禁（UpdateBlocked），如需定时解封可后续扩展。
func (a *forumAPIImpl) BanUser(ctx context.Context, userID uint, reason string, durationSec int) error {
	uid := uint(userID)
	if err := a.userRepo.UpdateBlocked(ctx, uid, true); err != nil {
		return err
	}
	// 发送封禁通知
	until := time.Now().Add(time.Duration(durationSec) * time.Second)
	notif := &do.Notification{
		UserID:     uid,
		Type:       do.NotifyBan,
		Content:    fmt.Sprintf("您的账号因 [%s] 已被封禁至 %s", reason, until.Format("2006-01-02 15:04")),
		TargetType: "user",
	}
	return a.notifRepo.Create(notif)
}
