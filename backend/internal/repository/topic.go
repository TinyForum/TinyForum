package repository

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type TopicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) *TopicRepository {
	return &TopicRepository{db: db}
}

func (r *TopicRepository) Create(topic *model.Topic) error {
	return r.db.Create(topic).Error
}

func (r *TopicRepository) Update(topic *model.Topic) error {
	return r.db.Save(topic).Error
}

func (r *TopicRepository) Delete(id uint) error {
	return r.db.Delete(&model.Topic{}, id).Error
}

func (r *TopicRepository) FindByID(id uint) (*model.Topic, error) {
	var topic model.Topic
	err := r.db.Preload("Creator").First(&topic, id).Error
	return &topic, err
}

func (r *TopicRepository) List(limit, offset int) ([]model.Topic, int64, error) {
	var topics []model.Topic
	var total int64

	query := r.db.Model(&model.Topic{})
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("follower_count DESC, post_count DESC").
		Find(&topics).Error
	return topics, total, err
}

func (r *TopicRepository) GetByCreator(creatorID uint, limit, offset int) ([]model.Topic, int64, error) {
	var topics []model.Topic
	var total int64

	query := r.db.Model(&model.Topic{}).Where("creator_id = ?", creatorID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&topics).Error
	return topics, total, err
}

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

// TopicPost methods
func (r *TopicRepository) AddPost(topicPost *model.TopicPost) error {
	// Check if already exists
	var existing model.TopicPost
	err := r.db.Where("topic_id = ? AND post_id = ?", topicPost.TopicID, topicPost.PostID).
		First(&existing).Error

	if err == nil {
		return nil // Already exists
	}

	return r.db.Create(topicPost).Error
}

func (r *TopicRepository) RemovePost(topicID, postID uint) error {
	return r.db.Where("topic_id = ? AND post_id = ?", topicID, postID).
		Delete(&model.TopicPost{}).Error
}

func (r *TopicRepository) GetTopicPosts(topicID uint, limit, offset int) ([]model.TopicPost, int64, error) {
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

func (r *TopicRepository) UpdatePostOrder(topicID, postID uint, sortOrder int) error {
	return r.db.Model(&model.TopicPost{}).
		Where("topic_id = ? AND post_id = ?", topicID, postID).
		Update("sort_order", sortOrder).Error
}

// TopicFollow methods
func (r *TopicRepository) Follow(follow *model.TopicFollow) error {
	// Check if already exists
	var existing model.TopicFollow
	err := r.db.Where("user_id = ? AND topic_id = ?", follow.UserID, follow.TopicID).
		First(&existing).Error

	if err == nil {
		return nil // Already following
	}

	return r.db.Create(follow).Error
}

func (r *TopicRepository) Unfollow(userID, topicID uint) error {
	return r.db.Where("user_id = ? AND topic_id = ?", userID, topicID).
		Delete(&model.TopicFollow{}).Error
}

func (r *TopicRepository) IsFollowing(userID, topicID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.TopicFollow{}).
		Where("user_id = ? AND topic_id = ?", userID, topicID).
		Count(&count).Error
	return count > 0, err
}

func (r *TopicRepository) GetFollowers(topicID uint, limit, offset int) ([]model.TopicFollow, int64, error) {
	var follows []model.TopicFollow
	var total int64

	query := r.db.Model(&model.TopicFollow{}).Where("topic_id = ?", topicID)
	query.Count(&total)

	err := query.Offset(offset).Limit(limit).
		Preload("User").
		Order("created_at DESC").
		Find(&follows).Error
	return follows, total, err
}
