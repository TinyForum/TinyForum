package topic

import (
	"tiny-forum/internal/repository"
	"tiny-forum/internal/service/notification"
)

type TopicService struct {
	topicRepo *repository.TopicRepository
	postRepo  repository.PostRepository
	userRepo  *repository.UserRepository
	notifSvc  *notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
}

func NewTopicService(
	topicRepo *repository.TopicRepository,
	postRepo repository.PostRepository,
	userRepo *repository.UserRepository,
	notifSvc *notification.NotificationService,
) *TopicService {
	return &TopicService{
		topicRepo: topicRepo,
		postRepo:  postRepo,
		userRepo:  userRepo,
		notifSvc:  notifSvc,
	}
}
