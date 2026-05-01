package user

import (
	"context"
	"tiny-forum/internal/model/dto"
)

func (s *userService) GetGlobalStatsCount(ctx context.Context, userID uint) (*dto.GlobalStatsCount, error) {
	return s.repo.GetGlobalStatsCount(ctx, userID)
}
