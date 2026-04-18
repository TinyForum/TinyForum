package user

import (
	"context"
	"fmt"
	"tiny-forum/internal/model"
)

// 返回排行榜原始数据（按积分降序，过滤被封禁用户）
func (s *UserService) GetLeaderboardData(ctx context.Context, limit int) ([]model.User, error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var users []model.User
	users, err := s.repo.GetTopScoreUsers(ctx, limit, true) // true 表示排除被封禁用户
	if err != nil {
		return nil, fmt.Errorf("查询排行榜失败: %w", err)
	}
	return users, nil
}
