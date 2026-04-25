package topic

import (
	postRepo "tiny-forum/internal/repository/post"
	topicRepo "tiny-forum/internal/repository/topic"
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"
)

type TopicService struct {
	topicRepo topicRepo.TopicRepository
	postRepo  postRepo.PostRepository
	userRepo  userRepo.UserRepository
	notifSvc  *notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
}

func NewTopicService(
	topicRepo topicRepo.TopicRepository,
	postRepo postRepo.PostRepository,
	userRepo userRepo.UserRepository,
	notifSvc *notification.NotificationService,
) *TopicService {
	return &TopicService{
		topicRepo: topicRepo,
		postRepo:  postRepo,
		userRepo:  userRepo,
		notifSvc:  notifSvc,
	}
}
