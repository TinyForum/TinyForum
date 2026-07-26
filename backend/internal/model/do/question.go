package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

type Question struct {
	common.BaseModel
	CreationID       uint     `gorm:"uniqueIndex;not null" json:"creations_id"`
	Creation         Creation `gorm:"foreignKey:CreationID;references:ID" json:"creation,omitempty"`
	AcceptedAnswerID *uint    `json:"accepted_answer_id"`
	RewardScore      int      `gorm:"default:0" json:"reward_score"`
	AnswerCount      int      `gorm:"default:0" json:"answer_count"`
	ViewCount        int      `gorm:"default:0" json:"view_count"`
	AcceptedAnswer   Comment  `gorm:"foreignKey:AcceptedAnswerID" json:"accepted_answer,omitempty"`
}

// CreateQuestionInput 创建问答输入

type QuestionResponse struct {
	ID               uint      `json:"id"`
	CreationsID      uint      `json:"creations_id"`
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
	common.BaseModel
	Title       string `json:"title"`
	Summary     string `json:"summary"`
	BoardID     uint   `json:"board_id"`
	AuthorID    uint   `json:"author_id"`
	RewardScore int    `json:"reward_score"`
	AnswerCount int    `json:"answer_count"`
	ViewCount   int    `gorm:"default:0" json:"view_count"`
}
