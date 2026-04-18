package user

import (
	"context"
	"fmt"
	"tiny-forum/internal/dto"
)

// 返回排行榜原始数据（按积分降序，过滤被封禁用户）
func (s *UserService) GetSimpleLeaderboardData(ctx context.Context, limit int) ([]dto.LeaderboardUserSimple, error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var users []dto.LeaderboardUserSimple
	users, err := s.repo.GetTopScoreUsersSimple(ctx, limit, true)
	if err != nil {
		return nil, fmt.Errorf("查询排行榜失败: %w", err)
	}
	return users, nil
}

func (s *UserService) GetDetailLeaderboardData(ctx context.Context, limit int) ([]dto.LeaderboardUserDetail, error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	var users []dto.LeaderboardUserDetail
	users, err := s.repo.GetTopScoreUsersDetail(ctx, limit, true)
	if err != nil {
		return nil, fmt.Errorf("查询排行榜失败: %w", err)
	}
	return users, nil
}
