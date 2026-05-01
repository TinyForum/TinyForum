package do

import (
	"time"
)

type Question struct {
	BaseModel
	PostID           uint    `gorm:"uniqueIndex;not null" json:"post_id"`
	AcceptedAnswerID *uint   `json:"accepted_answer_id"`
	RewardScore      int     `gorm:"default:0" json:"reward_score"`
	AnswerCount      int     `gorm:"default:0" json:"answer_count"`
	ViewCount        int     `gorm:"default:0" json:"view_count"`
	Post             Post    `gorm:"foreignKey:PostID" json:"post,omitempty"`
	AcceptedAnswer   Comment `gorm:"foreignKey:AcceptedAnswerID" json:"accepted_answer,omitempty"`
}

type AnswerVote struct {
	BaseModel
	UserID    uint   `gorm:"uniqueIndex:idx_user_answer;not null" json:"user_id"`
	CommentID uint   `gorm:"uniqueIndex:idx_user_answer;not null" json:"comment_id"`
	VoteType  string `gorm:"type:varchar(10)" json:"vote_type"` // up/down
}

// CreateQuestionInput 创建问答输入

type QuestionResponse struct {
	ID               uint      `json:"id"`
	PostID           uint      `json:"post_id"`
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	Summary          string    `json:"summary"`
	Cover            string    `json:"cover"`
	BoardID          uint      `json:"board_id"`
	AuthorID         uint      `json:"author_id"`
	RewardScore      int       `json:"reward_score"`
	AnswerCount      int       `json:"answer_count"`
	AcceptedAnswerID *uint     `json:"accepted_answer_id"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type QuestionListResponse struct {
	BaseModel
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	BoardID     uint   `json:"board_id"`
	AuthorID    uint   `json:"author_id"`
	RewardScore int    `json:"reward_score"`
	AnswerCount int    `json:"answer_count"`
	ViewCount   int    `gorm:"default:0" json:"view_count"`
}
