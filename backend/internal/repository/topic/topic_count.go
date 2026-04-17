package topic

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

func (r *TopicRepository) IncrementPostCount(topicID uint) error {
	return r.db.Model(&model.Topic{}).Where("id = ?", topicID).
		UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}

func (r *TopicRepository) DecrementPostCount(topicID uint) error {
	return r.db.Model(&model.Topic{}).Where("id = ?", topicID).
		UpdateColumn("post_count", gorm.Expr("post_count - 1")).Error
}

func (r *TopicRepository) IncrementFollowerCount(topicID uint) error {
	return r.db.Model(&model.Topic{}).Where("id = ?", topicID).
		UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error
}

func (r *TopicRepository) DecrementFollowerCount(topicID uint) error {
	return r.db.Model(&model.Topic{}).Where("id = ?", topicID).
		UpdateColumn("follower_count", gorm.Expr("follower_count - 1")).Error
}
