package request

type LeaderboardRequest struct {
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
	// Fields string `form:"fields"`
}

// SetUserRoleRequest 设置用户角色请求
type SetUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user member moderator reviewer bot admin super_admin system_maintainer"`
}

type GetUserPostsRequest struct {
	PageRequest
	Keyword          string `form:"keyword"`
	Status           string `form:"status"`            // 用户感知状态 (draft, published, archived)
	ModerationStatus string `form:"moderation_status"` // 风控状态审核结果 (normal, pending, rejected)
	Tag              string `form:"tag"`               // 标签名称（注意不是 ID）
	BoardName        string `form:"board_name"`
}

// GetUserCommentsRequest 我的评论查询参数
type GetUserCommentsRequest struct {
	PageRequest
}

// ListUserViolationRequest 用户违规记录查询参数
type ListUserViolationRequest struct {
	PageRequest
	Status  string `form:"status"`  // 违规处理状态 (pending, confirmed, rejected)
	Keyword string `form:"keyword"` // 关键词，模糊匹配违规原因
}

// AppealViolationRequest 违规申诉请求
type AppealViolationRequest struct {
	Reason string `json:"reason" binding:"required,min=1,max=500"` // 申诉理由
}
type AdminSetScoreRequest struct {
	Operation string `json:"operation" binding:"required,oneof=set add subtract"`
	Score     int    `json:"score" binding:"required,gte=0,lte=999999"`
	Reason    string `json:"reason" binding:"required,max=200"`
}
