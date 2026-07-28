package vo

import (
	"time"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
)

// QuestionVO 问答信息脱敏视图
type QuestionVO struct {
	ID               uint      `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	PostID           uint      `json:"post_id"`                      // 关联的帖子 ID
	AcceptedAnswerID *uint     `json:"accepted_answer_id,omitempty"` // 采纳的回答（评论）ID
	RewardScore      int       `json:"reward_score"`                 // 悬赏积分
	AnswerCount      int       `json:"answer_count"`                 // 回答数量
	ViewCount        int       `json:"view_count"`                   // 浏览次数
}

// AnswerVoteVO 回答投票脱敏视图
type AnswerVoteVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	CommentID uint      `json:"comment_id"`
	VoteType  string    `json:"vote_type"` // "up" 或 "down"
}

// QuestionSimpleData 问题精简数据
type QuestionSimpleDataVO struct {
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

// QuestionSimpleResponse 问题精简列表响应
type QuestionSimpleVO struct {
	ID               uint              `json:"id"`
	Title            string            `json:"title"`
	Summary          string            `json:"summary"`
	ViewCount        int               `json:"view_count"`
	RewardScore      int               `json:"reward_score"`
	AnswerCount      int               `json:"answer_count"`
	AcceptedAnswerID *uint             `json:"accepted_answer_id"`
	Author           *dto.SimpleAuthor `json:"author"`
	Tags             []dto.SimpleTag   `json:"tags"`
	CreatedAt        time.Time         `json:"created_at"`
}

// AcceptAnswerVO 采纳答案
type AcceptAnswerVO struct {
	Message QuestionActionType `json:"message"`
}

type QuestionActionType string

const (
	ActionAcceptedAnswer QuestionActionType = "已接受答案"
)

// CreateQuestionInput 创建问答输入

type QuestionDetailVO struct {
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
	Tags             []do.Tag  `json:"tags"`
	Author           UserVO    `json:"author"`
	ViewCount        int       `json:"view_count"`
	LikeCount        int       `json:"like_count"`
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
