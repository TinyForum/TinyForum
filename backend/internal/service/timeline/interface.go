package timeline

import (
	"tiny-forum/internal/model/po"
	commentRepo "tiny-forum/internal/repository/comment"
	postRepo "tiny-forum/internal/repository/post"
	timelineRepo "tiny-forum/internal/repository/timeline"
	userRepo "tiny-forum/internal/repository/user"
)

type TimelineService interface {
	CreateEvent(input CreateEventInput) error
	GetHomeTimeline(userID uint, page, pageSize int) ([]po.TimelineEvent, int64, error)
	GetFollowingTimeline(userID uint, page, pageSize int) ([]po.TimelineEvent, int64, error)
	Subscribe(subscriberID, targetUserID uint) error
	Unsubscribe(subscriberID, targetUserID uint) error
	GetSubscriptions(subscriberID uint) ([]po.TimelineSubscription, error)
	IsSubscribed(subscriberID, targetUserID uint) (bool, error)
}
type timelineService struct {
	timelineRepo timelineRepo.TimelineRepository
	userRepo     userRepo.UserRepository
	postRepo     postRepo.PostRepository
	commentRepo  commentRepo.CommentRepository
}

func NewTimelineService(
	timelineRepo timelineRepo.TimelineRepository,
	userRepo userRepo.UserRepository,
	postRepo postRepo.PostRepository,
	commentRepo commentRepo.CommentRepository,
) TimelineService {
	return &timelineService{
		timelineRepo: timelineRepo,
		userRepo:     userRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
	}
}
