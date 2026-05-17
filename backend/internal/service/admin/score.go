package admin

import (
	"context"
	"tiny-forum/internal/model/vo"
	apperrors "tiny-forum/pkg/errors"
)

// 列出所用用户用户积分
func (s *adminService) ListUsersScore(ctx context.Context) ([]vo.UserScoreVO, error) {
	users, err := s.userSvc.ListUsersScore() //GetEveryoneUsersScore()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// 获取单个用户积分
func (s *adminService) GetUserScore(ctx context.Context, userID uint) (*vo.UserScoreVO, error) {
	score, err := s.userSvc.GetScoreById(userID) //GetEveryoneUsersScore()
	if err != nil {
		return nil, apperrors.ErrInternalError
	}
	scoreVO := &vo.UserScoreVO{
		ID:    userID,
		Score: score,
	}
	return scoreVO, nil
}
