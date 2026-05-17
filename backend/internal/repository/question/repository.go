package question

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/vo"

	"gorm.io/gorm"
)

type QuestionRepository interface {
	// CRUD
	Create(question *do.Question) error
	Update(question *do.Question) error
	FindByID(id uint) (*do.Question, error)
	FindByPostID(postID uint) (*do.Question, error)
	IncrementAnswerCount(postID uint) error
	SetAcceptedAnswer(postID, commentID uint) error
	UpdateCommentVoteCount(commentID uint, voteCount int) error
	UpdateAnswerCount(questionID uint) error
	UpdateAcceptedAnswer(questionID uint, answerID uint) error

	// transaction
	CreateWithTransaction(userID uint, input dto.CreateQuestionRequest) (*do.QuestionResponse, error)

	// about vote
	CreateAnswerVote(vote *do.AnswerVote) error
	UpdateAnswerVote(vote *do.AnswerVote) error
	DeleteAnswerVote(userID, commentID uint) error
	FindAnswerVote(userID, commentID uint) (*do.AnswerVote, error)
	GetAnswerVoteCount(commentID uint) (int, error)

	// query
	FindSimple(pageSize, offset int, boardID *uint) ([]do.QuestionListResponse, int64, error)
	FindSimpleQuestions(pageSize, offset int, boardID *uint, filter, sort, keyword string) ([]vo.QuestionSimpleDataVO, int64, error)
	FindQuestionSimpleByID(questionID uint) (*vo.QuestionSimpleDataVO, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}
