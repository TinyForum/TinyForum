package timeline

import (
	"tiny-forum/internal/model"
)

// Subscribe 关注用户（如果已存在则重新激活）
func (r *timelineRepository) Subscribe(sub *model.TimelineSubscription) error {
	var existing model.TimelineSubscription
	err := r.db.Where("subscriber_id = ? AND target_user_id = ?", sub.SubscriberID, sub.TargetUserID).
		First(&existing).Error

	if err == nil {
		// 已存在，重新激活
		return r.db.Model(&existing).Update("is_active", true).Error
	}
	return r.db.Create(sub).Error
}

// Unsubscribe 取消关注（软删除，设置 is_active = false）
func (r *timelineRepository) Unsubscribe(subscriberID, targetUserID uint) error {
	return r.db.Model(&model.TimelineSubscription{}).
		Where("subscriber_id = ? AND target_user_id = ?", subscriberID, targetUserID).
		Update("is_active", false).Error
}

// GetSubscriptions 获取用户的所有有效关注
func (r *timelineRepository) GetSubscriptions(subscriberID uint) ([]model.TimelineSubscription, error) {
	var subs []model.TimelineSubscription
	err := r.db.Where("subscriber_id = ? AND is_active = ?", subscriberID, true).
		Find(&subs).Error
	return subs, err
}

// IsSubscribed 检查用户是否已关注指定用户
func (r *timelineRepository) IsSubscribed(subscriberID, targetUserID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.TimelineSubscription{}).
		Where("subscriber_id = ? AND target_user_id = ? AND is_active = ?",
			subscriberID, targetUserID, true).
		Count(&count).Error
	return count > 0, err
}
