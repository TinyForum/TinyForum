package user

import (
	"errors"
	"fmt"
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

// GetAllUsersWithScore 获取所有用户积分（用于管理员）
func (s *userService) GetAllUsersWithScore() ([]UserScoreResponse, error) {
	users, err := s.repo.GetEveryoneUsersScore()
	if err != nil {
		return nil, err
	}

	var result []UserScoreResponse
	for _, user := range users {
		basicInfo, err := s.repo.GetUserBasicInfoById(user.ID)
		if err != nil {
			continue
		}
		result = append(result, UserScoreResponse{
			ID:       user.ID,
			Username: basicInfo.Username,
			Avatar:   basicInfo.Avatar,
			Score:    user.Score,
		})
	}
	return result, nil
}
