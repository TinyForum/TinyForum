package vo

import (
	"time"
	"tiny-forum/internal/model/do"
)

// UserVO 用户脱敏视图（对外暴露）
type UserVO struct {
	ID          uint        `json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Username    string      `json:"username"`
	Avatar      string      `json:"avatar"`
	Bio         string      `json:"bio"`
	Role        do.UserRole `json:"role"`
	Score       int         `json:"score"`
	IsActive    bool        `json:"is_active"`
	IsBlocked   bool        `json:"is_blocked"`
	LastLogin   *time.Time  `json:"last_login,omitempty"`
	InvitedByID *uint       `json:"invited_by_id,omitempty"`
}

type UserPrivateVO struct {
	ID          uint        `json:"id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Username    string      `json:"username"`
	Avatar      string      `json:"avatar"`
	Bio         string      `json:"bio"`
	Role        do.UserRole `json:"role"`
	Score       int         `json:"score"`
	IsActive    bool        `json:"is_active"`
	IsBlocked   bool        `json:"is_blocked"`
	LastLogin   *time.Time  `json:"last_login,omitempty"`
	InvitedByID *uint       `json:"invited_by_id,omitempty"`
	Email       string      `json:"email"`
}

type UserPosts struct {
	ID               uint                `json:"id"`                                          // 帖子ID
	CreatedAt        time.Time           `json:"created_at"`                                  // 创建时间
	UpdatedAt        time.Time           `json:"updated_at"`                                  // 更新时间
	Title            string              `gorm:"not null;size:200" json:"title"`              // 标题
	Summary          string              `gorm:"size:500" json:"summary"`                     // 摘要
	Cover            string              `gorm:"size:500" json:"cover"`                       // 封面
	Type             do.PostType         `gorm:"type:varchar(20);default:'post'" json:"type"` // 帖子类型
	PostStatus       do.PostStatus       `json:"post_status"`                                 // 文章状态：主动状态
	ModerationStatus do.ModerationStatus `json:"moderation_status"`                           // 审核状态：被动状态
	ViewCount        int                 `json:"view_count"`                                  // 浏览数
	LikeCount        int                 `json:"likes_count"`                                 // 点赞数
	CommentCount     int64               `json:"comment_count"`                               // 新增评论数
	PinTop           bool                `json:"pin_top"`                                     // 用户主页置顶
	Tags             []string            `json:"tags"`                                        // 标签列表
	BoardName        string              `gorm:"index" json:"board_name"`                     // 所属板块
	PinInBoard       bool                `gorm:"default:false" json:"pin_in_board"`           // 板块置顶
}

// 不包含手机号、邮箱、IP 等
type UserPublicVO struct {
	ID     uint   `json:"id"`       // 用户ID
	Name   string `json:"nickname"` // 用户昵称
	Avatar string `json:"avatar"`   // 用户头像
}

type UserProfileVO struct {
	*do.User
	FollowerCount  int64 `json:"follower_count"`
	FollowingCount int64 `json:"following_count"`
	IsFollowing    bool  `json:"is_following"`
}
