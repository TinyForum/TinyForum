package user

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model/vo"
)

// GetScoreById 获取用户积分
func (s *userService) GetScoreById(userID uint) (int, error) {
	return s.repo.GetScoreById(userID)
}

// SetScoreById 设置用户积分
func (s *userService) SetScoreById(userID uint, score int) error {
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}
	err := s.repo.SetScoreById(userID, score)
	if err != nil {
		return fmt.Errorf("设置积分失败: %w", err)
	}
	go s.onScoreChanged(userID, score)
	return nil
}

// onScoreChanged 积分变更后的回调
func (s *userService) onScoreChanged(userID uint, newScore int) error {
	// 可扩展：发送通知、更新缓存等
	return nil
}

// 列出所有用户积分
func (s *userService) ListUsersScore() ([]vo.UserScoreVO, error) {
	scoreVO, err := s.repo.ListUsersScore()
	if err != nil {
		return nil, err
	}

	return scoreVO, nil
}
