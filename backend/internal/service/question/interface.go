package question

import (
	"context"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	postRepo "tiny-forum/internal/repository/article"
	commentRepo "tiny-forum/internal/repository/comment"
	questionRepo "tiny-forum/internal/repository/question"
	tagRepo "tiny-forum/internal/repository/tag"
	"tiny-forum/internal/repository/transaction"
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"
)

type QuestionService interface {
	AcceptAnswer(ctx context.Context, questionID, answerID uint, userID uint) error
	VoteAnswer(userID uint, input request.VoteAnswerRequest) (*vo.VoteAnswerVO, error)
	GetAnswerVoteStatus(userID, commentID uint) (map[string]any, error)
	GetAnswersList(questionID uint, page, pageSize int) ([]do.Answer, int64, error)
	// crud
	CreateQuestion(userID uint, input dto.CreateQuestionRequest) (*vo.QuestionDetailVO, error)
	GetQuestionDetail(questionID uint) (*vo.QuestionDetailVO, error)
	GetQuestionsList(page, pageSize int, unanswered bool) ([]do.Question, int64, error)
	GetQuestionByID(questionID uint) (*do.Question, error)
	// simple
	GetQuestionSimpleList(pageSize, offset int, boardID *uint, filter, sort, keyword string) ([]vo.QuestionSimpleVO, int64, error)
	GetQuestionSimpleByID(questionID uint) (*vo.QuestionSimpleVO, error)
	// GetAnswersByQuestion(questionID uint, page, pageSize int, sortBy string) ([]do.Answer, int64, error)
}

type questionService struct {
	questionRepo questionRepo.QuestionRepository
	postRepo     postRepo.ArticleRepository
	commentRepo  commentRepo.CommentRepository
	userRepo     userRepo.UserRepository
	notifSvc     notification.NotificationService
	tagRepo      tagRepo.TagRepository
	txManager    transaction.TransactionManager
}

func NewQuestionService(
	questionRepo questionRepo.QuestionRepository,
	postRepo postRepo.ArticleRepository,
	commentRepo commentRepo.CommentRepository,
	userRepo userRepo.UserRepository,
	notifSvc notification.NotificationService,
	tagRepo tagRepo.TagRepository,
	txManager transaction.TransactionManager,
) QuestionService {
	return &questionService{
		questionRepo: questionRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
		userRepo:     userRepo,
		notifSvc:     notifSvc,
		tagRepo:      tagRepo,
		txManager:    txManager,
	}
}
