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
