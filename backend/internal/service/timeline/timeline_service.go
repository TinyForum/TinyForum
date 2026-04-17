package timeline

import (
	"tiny-forum/internal/repository"
)

type TimelineService struct {
	timelineRepo *repository.TimelineRepository
	userRepo     *repository.UserRepository
	postRepo     repository.PostRepository
	commentRepo  *repository.CommentRepository
}

func NewTimelineService(
	timelineRepo *repository.TimelineRepository,
	userRepo *repository.UserRepository,
	postRepo repository.PostRepository,
	commentRepo *repository.CommentRepository,
) *TimelineService {
	return &TimelineService{
		timelineRepo: timelineRepo,
		userRepo:     userRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
	}
}
