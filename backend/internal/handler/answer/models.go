package answer

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
	UserVote  int `json:"user_vote"`  // 0:未投票, 1:赞同, -1:反对
	UpCount   int `json:"up_count"`   // 赞同数
	DownCount int `json:"down_count"` // 反对数
	// Total     int `json:"total"`     // 总投票数（可选）
}
