package question

import (
	"tiny-forum/internal/repository"
	"tiny-forum/internal/service/notification"

	"gorm.io/gorm"
)

type QuestionService struct {
	questionRepo repository.QuestionRepository
	postRepo     repository.PostRepository
	commentRepo  *repository.CommentRepository
	userRepo     *repository.UserRepository
	notifSvc     *notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
	db           *gorm.DB
	tagRepo      *repository.TagRepository
}

func NewQuestionService(
	questionRepo repository.QuestionRepository,
	postRepo repository.PostRepository,
	commentRepo *repository.CommentRepository,
	userRepo *repository.UserRepository,
	notifSvc *notification.NotificationService,
	tagRepo *repository.TagRepository,
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
