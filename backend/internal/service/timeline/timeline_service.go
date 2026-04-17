package timeline

import (
	commentRepo "tiny-forum/internal/repository/comment"
	postRepo "tiny-forum/internal/repository/post"
	timelineRepo "tiny-forum/internal/repository/timeline"
	userRepo "tiny-forum/internal/repository/user"
)

type TimelineService struct {
	timelineRepo *timelineRepo.TimelineRepository
	userRepo     *userRepo.UserRepository
	postRepo     postRepo.PostRepository
	commentRepo  *commentRepo.CommentRepository
}

func NewTimelineService(
	timelineRepo *timelineRepo.TimelineRepository,
	userRepo *userRepo.UserRepository,
	postRepo postRepo.PostRepository,
	commentRepo *commentRepo.CommentRepository,
) *TimelineService {
	return &TimelineService{
		timelineRepo: timelineRepo,
		userRepo:     userRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
	}
}
