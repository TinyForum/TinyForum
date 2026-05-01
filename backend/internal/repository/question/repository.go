package question

import (
	"time"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"

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
	FindSimpleQuestions(pageSize, offset int, boardID *uint, filter, sort, keyword string) ([]QuestionSimpleData, int64, error)
	FindQuestionSimpleByID(questionID uint) (*QuestionSimpleData, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

// QuestionSimpleData 问题精简数据（用于列表查询）
type QuestionSimpleData struct {
	ID               uint      `gorm:"column:id"`
	PostID           uint      `gorm:"column:post_id"`
	Title            string    `gorm:"column:title"`
	Summary          string    `gorm:"column:summary"`
	ViewCount        int       `gorm:"column:view_count"`
	BoardID          uint      `gorm:"column:board_id"`
	AuthorID         uint      `gorm:"column:author_id"`
	RewardScore      int       `gorm:"column:reward_score"`
	AnswerCount      int       `gorm:"column:answer_count"`
	AcceptedAnswerID *uint     `gorm:"column:accepted_answer_id"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}
