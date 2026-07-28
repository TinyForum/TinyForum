package vo

import (
	"time"
	"tiny-forum/internal/model/do"
)

// VoteVO 投票记录脱敏视图（对外暴露）
type VoteVO struct {
	ID        uint      `json:"id"`
	CommentID uint      `json:"comment_id"`
	Value     int       `json:"value"` // 1: 赞同, -1: 反对, 0: 取消
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type VoteAnswerVO struct {
	VoteType  *do.AnswerVoteType `json:"vote_type"`
	VoteCount int                `json:"vote_count"`
	Action    string             `json:"action"`
}

type VoteResponseVO struct {
	Message   string            `json:"message"`
	VoteCount int               `json:"vote_count"`
	UserVote  do.AnswerVoteType `json:"user_vote"` // 值为 "up", "down" 或 null
}
