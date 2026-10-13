package request

// ReviewTaskRequest 审核任务处理请求
type ReviewTaskRequest struct {
	Note string `json:"note" binding:"max=500"` // 审核备注/拒绝原因
}
