package question

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	commentRepo "tiny-forum/internal/repository/comment"
	postRepo "tiny-forum/internal/repository/post"
	questionRepo "tiny-forum/internal/repository/question"
	tagRepo "tiny-forum/internal/repository/tag"
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"

	"gorm.io/gorm"
)

type QuestionService interface {
	AcceptAnswer(postID, commentID uint, userID uint) error
	VoteAnswer(userID uint, input request.VoteAnswerRequest) (*vo.VoteAnswerVO, error)
	GetAnswerVoteStatus(userID, commentID uint) (map[string]interface{}, error)
	GetQuestionWithAnswers(postID uint, page, pageSize int) (*do.Question, []do.Comment, int64, error)
	// crud
	CreateQuestion(userID uint, input dto.CreateQuestionRequest) (*do.QuestionResponse, error)
	GetQuestionDetail(questionID uint) (*do.QuestionResponse, error)
	GetQuestionsList(page, pageSize int, unanswered bool) ([]do.Post, int64, error)
	GetQuestionByID(questionID uint) (*do.Question, error)
	// simple
	GetQuestionSimpleList(pageSize, offset int, boardID *uint, filter, sort, keyword string) ([]vo.QuestionSimpleVO, int64, error)
	GetQuestionSimpleByID(questionID uint) (*vo.QuestionSimpleVO, error)
}

type questionService struct {
	questionRepo questionRepo.QuestionRepository
	postRepo     postRepo.PostRepository
	commentRepo  commentRepo.CommentRepository
	userRepo     userRepo.UserRepository
	notifSvc     notification.NotificationService // 需导入 "tiny-forum/internal/service/notification"
	db           *gorm.DB
	tagRepo      tagRepo.TagRepository
}

func NewQuestionService(
	questionRepo questionRepo.QuestionRepository,
	postRepo postRepo.PostRepository,
	commentRepo commentRepo.CommentRepository,
	userRepo userRepo.UserRepository,
	notifSvc notification.NotificationService,
	tagRepo tagRepo.TagRepository,
) QuestionService {
	return &questionService{
		questionRepo: questionRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
		userRepo:     userRepo,
		notifSvc:     notifSvc,
		tagRepo:      tagRepo,
	}
}
