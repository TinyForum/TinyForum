package user

import (
	"context"
	"fmt"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
	"tiny-forum/pkg/fields"
)

// GetLeaderboard 获取排行榜
func (s *UserService) GetLeaderboard(ctx context.Context, limit int, fieldsParam string) ([]LeaderboardItem, error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	selectedFields := fields.Filter(
		fieldsParam,
		model.UserPublicFields,
		model.UserDefaultFields,
	)

	query := repository.TopUsersQuery{
		Limit:          limit,
		ExcludeBlocked: true,
		Fields:         selectedFields,
	}

	users, err := s.repo.GetTopUsers(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("查询排行榜失败: %w", err)
	}

	items := make([]LeaderboardItem, len(users))
	for i, u := range users {
		items[i] = LeaderboardItem{
			ID:       u.ID,
			Username: u.Username,
			Avatar:   u.Avatar,
			Score:    u.Score,
			Rank:     i + 1,
		}
	}
	return items, nil
}
