package vo

import "time"

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
