package post

import (
	boardRepo "tiny-forum/internal/repository/board"
	postRepo "tiny-forum/internal/repository/post"
	tagRepo "tiny-forum/internal/repository/tag"
	userRepo "tiny-forum/internal/repository/user"

	"tiny-forum/internal/service/notification"
)

type PostService struct {
	postRepo  postRepo.PostRepository
	tagRepo   *tagRepo.TagRepository
	boardRepo *boardRepo.BoardRepository
	userRepo  *userRepo.UserRepository
	notifSvc  *notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
}

func NewPostService(
	postRepo postRepo.PostRepository,
	tagRepo *tagRepo.TagRepository,
	userRepo *userRepo.UserRepository,
	boardRepo *boardRepo.BoardRepository,
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
