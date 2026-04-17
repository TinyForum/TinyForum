package comment

import (
	"tiny-forum/internal/repository"
	"tiny-forum/internal/service/notification"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    *repository.UserRepository
	notifSvc    *notification.NotificationService // 注意：需正确导入包
	voteRepo    *repository.VoteRepository
}

func NewCommentService(
	commentRepo *repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo *repository.UserRepository,
	notifSvc *notification.NotificationService,
	voteRepo *repository.VoteRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		notifSvc:    notifSvc,
		voteRepo:    voteRepo,
	}
}
