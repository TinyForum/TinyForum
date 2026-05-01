package topic

import "tiny-forum/internal/model/do"

func (r *topicRepository) Follow(follow *do.TopicFollow) error {
	var existing do.TopicFollow
	err := r.db.Where("user_id = ? AND topic_id = ?", follow.UserID, follow.TopicID).
		First(&existing).Error

	if err == nil {
		return nil // Already following
	}
	return r.db.Create(follow).Error
}

func (r *topicRepository) Unfollow(userID, topicID uint) error {
	return r.db.Where("user_id = ? AND topic_id = ?", userID, topicID).
		Delete(&do.TopicFollow{}).Error
}

func (r *topicRepository) IsFollowing(userID, topicID uint) (bool, error) {
	var count int64
	err := r.db.Model(&do.TopicFollow{}).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Count(&count).Error
	return count > 0, err
}

func (r *topicRepository) GetFollowers(topicID uint, limit, offset int) ([]do.TopicFollow, int64, error) {
	var follows []do.TopicFollow
	var total int64

	query := r.db.Model(&do.TopicFollow{}).Where("topic_id = ?", topicID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("User").
		Order("created_at DESC").
		Find(&follows).Error
	return follows, total, err
}
