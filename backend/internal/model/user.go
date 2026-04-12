package model

import (
	"time"
)

// 定义 context key 常量
const (
	ContextUserID   = "user_id"
	ContextUsername = "username"
	ContextUserRole = "user_role"
)

// User 模型
type User struct {
	BaseModel
	Username  string     `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email     string     `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Password  string     `gorm:"not null" json:"-"`
	Avatar    string     `gorm:"size:500" json:"avatar"`
	Bio       string     `gorm:"size:500" json:"bio"`
	Role      UserRole   `gorm:"type:varchar(20);default:'user'" json:"role"`
	Score     int        `gorm:"default:0" json:"score"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	IsBlocked bool       `gorm:"default:false" json:"is_blocked"`
	LastLogin *time.Time `json:"last_login"`

	// Relations
	Posts     []Post    `gorm:"foreignKey:AuthorID" json:"-"`
	Comments  []Comment `gorm:"foreignKey:AuthorID" json:"-"`
	Followers []Follow  `gorm:"foreignKey:FollowingID" json:"-"`
	Following []Follow  `gorm:"foreignKey:FollowerID" json:"-"`
}

// 辅助方法：检查用户是否有某个权限
func (u *User) Can(perm Permission) bool {
	return HasPermission(u.Role, perm)
}

// 辅助方法：检查用户是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin || u.Role == RoleSuperAdmin
}

// 辅助方法：检查用户是否是版主
func (u *User) IsModerator() bool {
	return u.Role == RoleModerator || u.IsAdmin()
}
