package user

// LoginInput 登录请求
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResult 认证结果
// type AuthResult struct {
// 	// Deprecated: 停用
// 	Token string     `json:"token"` // token 目前没有用到，作为保留
// 	User  *vo.UserVO `json:"user"`
// }

// UserProfileResponse 用户资料响应（含关注统计）

// LeaderboardItem 排行榜条目
type LeaderboardItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}

// LoginResult 登录结果（可选）
type LoginResult struct {
	Token string    `json:"-"`
	User  *UserInfo `json:"user"`
}

// UserInfo 用户简要信息
type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
