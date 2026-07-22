package timeline

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	postRepo "tiny-forum/internal/repository/article"
	commentRepo "tiny-forum/internal/repository/comment"
	timelineRepo "tiny-forum/internal/repository/timeline"
	userRepo "tiny-forum/internal/repository/user"
)

type TimelineService interface {
	CreateEvent(input request.CreateEventRequest) error
	GetHomeTimeline(userID uint, page, pageSize int) ([]do.TimelineEvent, int64, error)
	GetFollowingTimeline(userID uint, page, pageSize int) ([]do.TimelineEvent, int64, error)
	Subscribe(subscriberID, targetUserID uint) error
	Unsubscribe(subscriberID, targetUserID uint) error
	GetSubscriptions(subscriberID uint) ([]do.TimelineSubscription, error)
	IsSubscribed(subscriberID, targetUserID uint) (bool, error)
}
type timelineService struct {
	timelineRepo timelineRepo.TimelineRepository
	userRepo     userRepo.UserRepository
	postRepo     postRepo.ArticleRepository
	commentRepo  commentRepo.CommentRepository
}

func NewTimelineService(
	timelineRepo timelineRepo.TimelineRepository,
	userRepo userRepo.UserRepository,
	postRepo postRepo.ArticleRepository,
	commentRepo commentRepo.CommentRepository,
) TimelineService {
	return &timelineService{
		timelineRepo: timelineRepo,
		userRepo:     userRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
	}
}
