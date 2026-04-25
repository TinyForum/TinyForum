package question

import (
	commentRepo "tiny-forum/internal/repository/comment"
	postRepo "tiny-forum/internal/repository/post"
	questionRepo "tiny-forum/internal/repository/question"
	tagRepo "tiny-forum/internal/repository/tag"
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"

	"gorm.io/gorm"
)

type QuestionService struct {
	questionRepo questionRepo.QuestionRepository
	postRepo     postRepo.PostRepository
	commentRepo  *commentRepo.CommentRepository
	userRepo     userRepo.UserRepository
	notifSvc     *notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
	db           *gorm.DB
	tagRepo      *tagRepo.TagRepository
}

func NewQuestionService(
	questionRepo questionRepo.QuestionRepository,
	postRepo postRepo.PostRepository,
	commentRepo *commentRepo.CommentRepository,
	userRepo userRepo.UserRepository,
	notifSvc *notification.NotificationService,
	tagRepo *tagRepo.TagRepository,
) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
		userRepo:     userRepo,
		notifSvc:     notifSvc,
		tagRepo:      tagRepo,
	}
}
