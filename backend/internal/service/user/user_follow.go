package user

import (
	"errors"
	"tiny-forum/internal/model"
)

// Follow 关注用户
func (s *userService) Follow(followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("不能关注自己")
	}
	if err := s.repo.Follow(followerID, followingID); err != nil {
		return err
	}
	following, _ := s.repo.FindByID(followingID)
	if following != nil {
		s.notifSvc.Create(followingID, &followerID, model.NotifyFollow,
			following.Username+" 关注了你", nil, "")
	}
	return nil
}

// Unfollow 取消关注
func (s *userService) Unfollow(followerID, followingID uint) error {
	return s.repo.Unfollow(followerID, followingID)
}

// GetFollowers 获取粉丝列表
func (s *userService) GetFollowers(userID uint, page, pageSize int) ([]model.User, int64, error) {
	return s.repo.GetFollowers(userID, page, pageSize)
}

// GetFollowing 获取关注列表
func (s *userService) GetFollowing(userID uint, page, pageSize int) ([]model.User, int64, error) {
	return s.repo.GetFollowing(userID, page, pageSize)
}
