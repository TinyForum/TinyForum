package topic

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

func (r *topicRepository) IncrementPostCount(topicID uint) error {
	return r.db.Model(&model.Topic{}).Where("id = ?", topicID).
		UpdateColumn("post_count", gorm.Expr("post_count + 1")).Error
}

func (r *topicRepository) DecrementPostCount(topicID uint) error {
	return r.db.Model(&model.Topic{}).Where("id = ?", topicID).
		UpdateColumn("post_count", gorm.Expr("post_count - 1")).Error
}

func (r *topicRepository) IncrementFollowerCount(topicID uint) error {
	return r.db.Model(&model.Topic{}).Where("id = ?", topicID).
		UpdateColumn("follower_count", gorm.Expr("follower_count + 1")).Error
}

func (r *topicRepository) DecrementFollowerCount(topicID uint) error {
	return r.db.Model(&model.Topic{}).Where("id = ?", topicID).
		UpdateColumn("follower_count", gorm.Expr("follower_count - 1")).Error
}
