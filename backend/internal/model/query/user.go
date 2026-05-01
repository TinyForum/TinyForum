package query

type LeaderboardRequest struct {
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
	// Fields string `form:"fields"`
}

// SetUserRoleRequest 设置用户角色请求
type SetUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user member moderator reviewer bot admin super_admin"`
}