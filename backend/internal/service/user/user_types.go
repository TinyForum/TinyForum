package user

import "tiny-forum/internal/model/po"

// RegisterInput 注册请求
type RegisterInput struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput 登录请求
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResult 认证结果
type AuthResult struct {
	Token string   `json:"token"`
	User  *po.User `json:"user"`
}

// UserProfileResponse 用户资料响应（含关注统计）
type UserProfileResponse struct {
	*po.User
	FollowerCount  int64 `json:"follower_count"`
	FollowingCount int64 `json:"following_count"`
	IsFollowing    bool  `json:"is_following"`
}

// LeaderboardItem 排行榜条目
type LeaderboardItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}

// UserScoreResponse 积分响应（用于列表）
type UserScoreResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar_url"`
	Score    int    `json:"score"`
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
