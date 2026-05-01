package topic

import (
	"tiny-forum/internal/model/po"
	postRepo "tiny-forum/internal/repository/post"
	topicRepo "tiny-forum/internal/repository/topic"
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"
)

type TopicService interface {
	Create(creatorID uint, input CreateTopicInput) (*po.Topic, error)
	Update(id uint, input CreateTopicInput) (*po.Topic, error)
	Delete(id uint, userID uint, isAdmin bool) error
	GetByID(id uint) (*po.Topic, error)
	List(page, pageSize int) ([]po.Topic, int64, error)
	GetByCreator(creatorID uint, page, pageSize int) ([]po.Topic, int64, error)
	// follow
	Follow(userID, topicID uint) error
	Unfollow(userID, topicID uint) error
	IsFollowing(userID, topicID uint) (bool, error)
	GetFollowers(topicID uint, page, pageSize int) ([]po.TopicFollow, int64, error)
	// post
	AddPostToTopic(input AddPostToTopicInput, userID uint) error
	RemovePostFromTopic(topicID, postID uint, userID uint) error
	GetTopicPosts(topicID uint, page, pageSize int) ([]po.TopicPost, int64, error)
}
type topicService struct {
	topicRepo topicRepo.TopicRepository
	postRepo  postRepo.PostRepository
	userRepo  userRepo.UserRepository
	notifSvc  notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
}

func NewTopicService(
	topicRepo topicRepo.TopicRepository,
	postRepo postRepo.PostRepository,
	userRepo userRepo.UserRepository,
	notifSvc notification.NotificationService,
) TopicService {
	return &topicService{
		topicRepo: topicRepo,
		postRepo:  postRepo,
		userRepo:  userRepo,
		notifSvc:  notifSvc,
	}
}
