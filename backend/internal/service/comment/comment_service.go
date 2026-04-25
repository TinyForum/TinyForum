package comment

import (
	commentRepo "tiny-forum/internal/repository/comment"
	postRepo "tiny-forum/internal/repository/post"
	userRepo "tiny-forum/internal/repository/user"
	voteRepo "tiny-forum/internal/repository/vote"
	"tiny-forum/internal/service/notification"
)

type CommentService struct {
	commentRepo commentRepo.CommentRepository
	postRepo    postRepo.PostRepository
	userRepo    userRepo.UserRepository
	notifSvc    notification.NotificationService
	voteRepo    voteRepo.VoteRepository
}

func NewCommentService(
	commentRepo commentRepo.CommentRepository,
	postRepo postRepo.PostRepository,
	userRepo userRepo.UserRepository,
	notifSvc notification.NotificationService,
	voteRepo voteRepo.VoteRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		notifSvc:    notifSvc,
		voteRepo:    voteRepo,
	}
}
