package model

import "time"

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleAdmin     UserRole = "admin"
	ModeratorUser UserRole = "moderator"
)

type User struct {
	// gorm.Model
	BaseModel
	Username  string     `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email     string     `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Password  string     `gorm:"not null" json:"-"`
	Avatar    string     `gorm:"size:500" json:"avatar"`
	Bio       string     `gorm:"size:500" json:"bio"`
	Role      UserRole   `gorm:"type:varchar(20);default:'user'" json:"role"`
	Score     int        `gorm:"default:0" json:"score"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`   // 激活，是否可以登录
	IsBlocked bool       `gorm:"default:false" json:"is_blocked"` // 封禁，是否可以发言
	LastLogin *time.Time `json:"last_login"`

	// Relations
	Posts     []Post    `gorm:"foreignKey:AuthorID" json:"-"`
	Comments  []Comment `gorm:"foreignKey:AuthorID" json:"-"`
	Followers []Follow  `gorm:"foreignKey:FollowingID" json:"-"`
	Following []Follow  `gorm:"foreignKey:FollowerID" json:"-"`
}
