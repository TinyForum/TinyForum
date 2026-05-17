package botapi

import (
	"context"
	"time"
	"tiny-forum/internal/infra/lua/sdk"
)

// ─── Stats ────────────────────────────────────────────────────────────────

func (a *forumAPIImpl) GetForumStats(ctx context.Context) (*sdk.StatsVO, error) {
	postCount, err := a.postRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	userCount, err := a.userRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	commentCount, err := a.commentRepo.Count(ctx)
	if err != nil {
		return nil, err
	}
	// 活跃用户（近 24h 内有注册的用户数作为近似值）
	yesterday := time.Now().Add(-24 * time.Hour)
	activeToday, _ := a.userRepo.CountActiveByDateRange(ctx, yesterday, time.Now())

	return &sdk.StatsVO{
		PostCount:    postCount,
		UserCount:    userCount,
		CommentCount: commentCount,
		ActiveToday:  activeToday,
	}, nil
}
