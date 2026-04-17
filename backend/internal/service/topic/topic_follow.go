package topic

import (
	"errors"

	"tiny-forum/internal/model"
)

// Follow 关注专题
func (s *TopicService) Follow(userID, topicID uint) error {
	_, err := s.topicRepo.FindByID(topicID)
	if err != nil {
		return errors.New("专题不存在")
	}
	isFollowing, _ := s.topicRepo.IsFollowing(userID, topicID)
	if isFollowing {
		return errors.New("已经关注过了")
	}
	follow := &model.TopicFollow{
		UserID:  userID,
		TopicID: topicID,
	}
	if err := s.topicRepo.Follow(follow); err != nil {
		return err
	}
	return s.topicRepo.IncrementFollowerCount(topicID)
}

// Unfollow 取消关注专题
func (s *TopicService) Unfollow(userID, topicID uint) error {
	_, err := s.topicRepo.FindByID(topicID)
	if err != nil {
		return errors.New("专题不存在")
	}
	if err := s.topicRepo.Unfollow(userID, topicID); err != nil {
		return err
	}
	return s.topicRepo.DecrementFollowerCount(topicID)
}

// IsFollowing 检查是否已关注专题
func (s *TopicService) IsFollowing(userID, topicID uint) (bool, error) {
	return s.topicRepo.IsFollowing(userID, topicID)
}

// GetFollowers 获取专题的关注者列表
func (s *TopicService) GetFollowers(topicID uint, page, pageSize int) ([]model.TopicFollow, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.topicRepo.GetFollowers(topicID, pageSize, offset)
}
