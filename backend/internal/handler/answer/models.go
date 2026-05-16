package answer

import "tiny-forum/internal/model/do"

// CreateAnswerRequest 提交回答的请求参数
type CreateAnswerRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"` // 回答内容
}

// VoteAnswerRequest 投票请求参数
type VoteAnswerRequest struct {
	VoteType string `json:"vote_type" binding:"required,oneof=up down" example:"up"` // up-赞同，down-反对
}

// VoteStatusResponse 投票状态响应
type VoteStatusResponse struct {
	UserVote  *do.AnswerVoteType `json:"user_vote,omitempty"` // 用户投票类型，未投票为 null
	UpCount   int                `json:"up_count"`            // 赞同数
	DownCount int                `json:"down_count"`          // 反对数
	Total     int                `json:"total"`               // 总投票数
}
