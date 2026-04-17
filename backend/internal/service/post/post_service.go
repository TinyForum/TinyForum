package post

import (
	"tiny-forum/internal/repository"
	"tiny-forum/internal/service/notification"
)

type PostService struct {
	postRepo  repository.PostRepository
	tagRepo   *repository.TagRepository
	boardRepo *repository.BoardRepository
	userRepo  *repository.UserRepository
	notifSvc  *notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
}

func NewPostService(
	postRepo repository.PostRepository,
	tagRepo *repository.TagRepository,
	userRepo *repository.UserRepository,
	boardRepo *repository.BoardRepository,
	notifSvc *notification.NotificationService,
) *PostService {
	return &PostService{
		postRepo:  postRepo,
		tagRepo:   tagRepo,
		userRepo:  userRepo,
		boardRepo: boardRepo,
		notifSvc:  notifSvc,
	}
}
