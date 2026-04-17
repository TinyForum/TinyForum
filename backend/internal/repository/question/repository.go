package question

import (
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

// QuestionRepository 问答数据访问接口
type QuestionRepository interface {
	// 基础 CRUD
	Create(question *model.Question) error
	Update(question *model.Question) error
	FindByID(id uint) (*model.Question, error)
	FindByPostID(postID uint) (*model.Question, error)
	IncrementAnswerCount(postID uint) error
	SetAcceptedAnswer(postID, commentID uint) error
	UpdateCommentVoteCount(commentID uint, voteCount int) error
	UpdateAnswerCount(questionID uint) error
	UpdateAcceptedAnswer(questionID uint, answerID uint) error

	// 事务
	CreateWithTransaction(userID uint, input model.CreateQuestionInput) (*model.QuestionResponse, error)

	// 投票
	CreateAnswerVote(vote *model.AnswerVote) error
	UpdateAnswerVote(vote *model.AnswerVote) error
	DeleteAnswerVote(userID, commentID uint) error
	FindAnswerVote(userID, commentID uint) (*model.AnswerVote, error)
	GetAnswerVoteCount(commentID uint) (int, error)

	// 查询
	FindSimple(pageSize, offset int, boardID *uint) ([]model.QuestionListResponse, int64, error)
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
