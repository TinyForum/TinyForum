package topic

import (
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

type TopicRepository interface {
	// count
	IncrementPostCount(topicID uint) error
	DecrementPostCount(topicID uint) error
	IncrementFollowerCount(topicID uint) error
	DecrementFollowerCount(topicID uint) error
	// curd
	Create(topic *po.Topic) error
	Update(topic *po.Topic) error
	Delete(id uint) error
	FindByID(id uint) (*po.Topic, error)
	List(limit, offset int) ([]po.Topic, int64, error)
	GetByCreator(creatorID uint, limit, offset int) ([]po.Topic, int64, error)
	// follow
	Follow(follow *po.TopicFollow) error
	Unfollow(userID, topicID uint) error
	IsFollowing(userID, topicID uint) (bool, error)
	GetFollowers(topicID uint, limit, offset int) ([]po.TopicFollow, int64, error)
	// post
	AddPost(topicPost *po.TopicPost) error
	RemovePost(topicID, postID uint) error
	GetTopicPosts(topicID uint, limit, offset int) ([]po.TopicPost, int64, error)
	UpdatePostOrder(topicID, postID uint, sortOrder int) error
}
type topicRepository struct {
	db *gorm.DB
}

func NewTopicRepository(db *gorm.DB) TopicRepository {
	return &topicRepository{db: db}
}
