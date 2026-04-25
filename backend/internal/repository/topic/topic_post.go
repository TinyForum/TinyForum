package topic

import (
	"tiny-forum/internal/model"
)

func (r *topicRepository) AddPost(topicPost *model.TopicPost) error {
	var existing model.TopicPost
	err := r.db.Where("topic_id = ? AND post_id = ?", topicPost.TopicID, topicPost.PostID).
		First(&existing).Error

	if err == nil {
		return nil // Already exists
	}
	return r.db.Create(topicPost).Error
}

func (r *topicRepository) RemovePost(topicID, postID uint) error {
	return r.db.Where("topic_id = ? AND post_id = ?", topicID, postID).
		Delete(&model.TopicPost{}).Error
}

func (r *topicRepository) GetTopicPosts(topicID uint, limit, offset int) ([]model.TopicPost, int64, error) {
	var topicPosts []model.TopicPost
	var total int64

	query := r.db.Model(&model.TopicPost{}).Where("topic_id = ?", topicID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Post").
		Preload("Post.Author").
		Order("sort_order ASC, created_at ASC").
		Find(&topicPosts).Error
	return topicPosts, total, err
}

func (r *topicRepository) UpdatePostOrder(topicID, postID uint, sortOrder int) error {
	return r.db.Model(&model.TopicPost{}).
		Where("topic_id = ? AND post_id = ?", topicID, postID).
		Update("sort_order", sortOrder).Error
}
