package topic

import "tiny-forum/internal/model/do"

// 添加帖子到主题
func (r *topicRepository) AddPost(topicPost *do.TopicCreation) error {
	var existing do.TopicCreation
	// 检查是否存在
	err := r.db.Where("topic_id = ? AND creation_id = ?", topicPost.TopicID, topicPost.CreationID).
		First(&existing).Error

	if err == nil {
		return nil
	}
	return r.db.Create(topicPost).Error
}

func (r *topicRepository) RemovePost(topicID, postID uint) error {
	return r.db.Unscoped().Where("topic_id = ? AND creation_id = ?", topicID, postID).
		Delete(&do.TopicCreation{}).Error
}

// 获取主题下的帖子
func (r *topicRepository) GetTopicPosts(topicID uint, limit, offset int) ([]do.TopicCreation, int64, error) {
	var topicPosts []do.TopicCreation
	var total int64

	query := r.db.Model(&do.TopicCreation{}).Where("topic_id = ?", topicID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("sort_order ASC, created_at ASC").
		Find(&topicPosts).Error
	return topicPosts, total, err
}

func (r *topicRepository) UpdatePostOrder(topicID, postID uint, sortOrder int) error {
	return r.db.Model(&do.TopicCreation{}).
		Where("topic_id = ? AND creation_id = ?", topicID, postID).
		Update("sort_order", sortOrder).Error
}
