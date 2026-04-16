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
	Username  string     `gorm:"uniqueIndex;not null;size:50" json:"username"` // 用户名
	Email     string     `gorm:"uniqueIndex;not null;size:100" json:"email"`   // 邮箱
	Password  string     `gorm:"not null" json:"-"`                            // 密码（hash）
	Avatar    string     `gorm:"size:500" json:"avatar"`                       // 用户的头像，目前为链接（暂不支持图片上传）
	Bio       string     `gorm:"size:500" json:"bio"`
	Role      UserRole   `gorm:"type:varchar(20);default:'user'" json:"role"` // 用户的角色，默认为普通用户
	Score     int        `gorm:"default:0" json:"score"`                      // 用户的积分
	IsActive  bool       `gorm:"default:true" json:"is_active"`               // 低优先级，用户主动行为，例如验证邮箱后可以处于激活状态非处罚性质
	IsBlocked bool       `gorm:"default:false" json:"is_blocked"`             // 优先级高于 IsActive，被动行为，一旦味 true 完全无法登录，处罚性质
	LastLogin *time.Time `json:"last_login"`                                  // 最后登陆时间

	IsTempPassword     bool       `gorm:"default:false" json:"-"` // 是否为临时密码
	TempPasswordExpire *time.Time `json:"-"`                      // 临时密码过期时间

	// 社交活动，可用于风控、审查
	Posts     []Post    `gorm:"foreignKey:AuthorID" json:"-"`    // 用户发布的帖子
	Comments  []Comment `gorm:"foreignKey:AuthorID" json:"-"`    // 用户发布的评论
	Followers []Follow  `gorm:"foreignKey:FollowingID" json:"-"` // 用户的粉丝
	Following []Follow  `gorm:"foreignKey:FollowerID" json:"-"`  // 用户关注的用户
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
