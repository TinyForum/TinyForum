package timeline

import "tiny-forum/internal/model/do"

// Subscribe 关注用户
func (s *timelineService) Subscribe(subscriberID, targetUserID uint) error {
	_, err := s.userRepo.FindByID(targetUserID)
	if err != nil {
		return err
	}
	sub := &do.TimelineSubscription{
		SubscriberID: subscriberID,
		TargetUserID: targetUserID,
		TargetType:   "user",
		IsActive:     true,
	}
	return s.timelineRepo.Subscribe(sub)
}

// Unsubscribe 取消关注用户
func (s *timelineService) Unsubscribe(subscriberID, targetUserID uint) error {
	return s.timelineRepo.Unsubscribe(subscriberID, targetUserID)
}

// GetSubscriptions 获取关注列表
func (s *timelineService) GetSubscriptions(subscriberID uint) ([]do.TimelineSubscription, error) {
	return s.timelineRepo.GetSubscriptions(subscriberID)
}

// IsSubscribed 检查是否已关注
func (s *timelineService) IsSubscribed(subscriberID, targetUserID uint) (bool, error) {
	return s.timelineRepo.IsSubscribed(subscriberID, targetUserID)
}
