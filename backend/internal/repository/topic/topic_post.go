package topic

import "tiny-forum/internal/model/do"

func (r *topicRepository) AddPost(topicPost *do.TopicPost) error {
	var existing do.TopicPost
	err := r.db.Where("topic_id = ? AND post_id = ?", topicPost.TopicID, topicPost.PostID).
		First(&existing).Error

	if err == nil {
		return nil // Already exists
	}
	return r.db.Create(topicPost).Error
}

func (r *topicRepository) RemovePost(topicID, postID uint) error {
	return r.db.Where("topic_id = ? AND post_id = ?", topicID, postID).
		Delete(&do.TopicPost{}).Error
}

func (r *topicRepository) GetTopicPosts(topicID uint, limit, offset int) ([]do.TopicPost, int64, error) {
	var topicPosts []do.TopicPost
	var total int64

	query := r.db.Model(&do.TopicPost{}).Where("topic_id = ?", topicID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("Post").
		Preload("Post.Author").
		Order("sort_order ASC, created_at ASC").
		Find(&topicPosts).Error
	return topicPosts, total, err
}

func (r *topicRepository) UpdatePostOrder(topicID, postID uint, sortOrder int) error {
	return r.db.Model(&do.TopicPost{}).
		Where("topic_id = ? AND post_id = ?", topicID, postID).
		Update("sort_order", sortOrder).Error
}
