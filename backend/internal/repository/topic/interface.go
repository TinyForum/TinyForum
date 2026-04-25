package topic

import (
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type TopicRepository interface {
	// count
	IncrementPostCount(topicID uint) error
	DecrementPostCount(topicID uint) error
	IncrementFollowerCount(topicID uint) error
	DecrementFollowerCount(topicID uint) error
	// curd
	Create(topic *model.Topic) error
	Update(topic *model.Topic) error
	Delete(id uint) error
	FindByID(id uint) (*model.Topic, error)
	List(limit, offset int) ([]model.Topic, int64, error)
	GetByCreator(creatorID uint, limit, offset int) ([]model.Topic, int64, error)
	// follow
	Follow(follow *model.TopicFollow) error
	Unfollow(userID, topicID uint) error
	IsFollowing(userID, topicID uint) (bool, error)
	GetFollowers(topicID uint, limit, offset int) ([]model.TopicFollow, int64, error)
	// post
	AddPost(topicPost *model.TopicPost) error
	RemovePost(topicID, postID uint) error
	GetTopicPosts(topicID uint, limit, offset int) ([]model.TopicPost, int64, error)
	UpdatePostOrder(topicID, postID uint, sortOrder int) error
}
type topicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return &topicRepository{db: db}
}
