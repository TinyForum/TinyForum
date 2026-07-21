package vo

import (
	"time"
	"tiny-forum/internal/model/do"
)

// 用户注册
type UserRegisterResultVO struct {
	Token string `json:"token"` // JWT 访问令牌
	// TokenType string `json:"token_type"` // "Bearer"
	// ExpiresIn int64  `json:"expires_in"` // 过期时间（秒）

	// 只返回前端渲染 UI 必要的字段
	User *UserRegisterVO `json:"user"`
}

type UserRegisterVO struct {
	ID        uint        `json:"id"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	AvatarUrl string      `json:"avatar_url"`
	Role      do.UserRole `json:"role"`
	Score     int         `json:"score"`
	CreatedAt time.Time   `json:"created_at"`
}

// 用户登录
type UserLoginResultVO struct {
	Token string `json:"token"` // JWT 访问令牌
	// TokenType string `json:"token_type"` // "Bearer"
	// ExpiresIn int64  `json:"expires_in"` // 过期时间（秒）

	// 只返回前端渲染 UI 必要的字段
	User *UserLoginVO `json:"user"`
}
type UserLoginVO struct {
	ID        uint        `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Username  string      `json:"username"`
	AvatarUrl string      `json:"avatar_url"`
	Bio       string      `json:"bio"`
	Role      do.UserRole `json:"role"`
	Score     int         `json:"score"`
	LastLogin *time.Time  `json:"last_login,omitempty"`
	Email     string      `json:"email"`
	// IsActive  bool        `json:"is_active"`
	// IsBlocked bool        `json:"is_blocked"`
	// LastLogin   *time.Time  `json:"last_login,omitempty"`
	InvitedByID *uint `json:"invited_by_id,omitempty"`
}
type AuthResultVO struct {
	Token          string          `json:"token"`
	User           *UserRegisterVO `json:"user"`
	DeletionStatus *DeletionStatus `json:"deletion_status,omitempty"`
}

type DeletionStatus struct {
	IsDeleted     bool       `json:"is_deleted"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	CanRestore    bool       `json:"can_restore"`
	RemainingDays int        `json:"remaining_days,omitempty"`
}

type DeleteAccountVO struct {
	IsDeleted bool `json:"is_deleted"`
}
