package model

import (
	"time"

	"gorm.io/gorm"
)

// ─── User ────────────────────────────────────────────────────────────────────

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	gorm.Model
	Username  string    `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null;size:100" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Avatar    string    `gorm:"size:500" json:"avatar"`
	Bio       string    `gorm:"size:500" json:"bio"`
	Role      UserRole  `gorm:"type:varchar(20);default:'user'" json:"role"`
	Score     int       `gorm:"default:0" json:"score"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time `json:"last_login"`

	// Relations
	Posts     []Post    `gorm:"foreignKey:AuthorID" json:"-"`
	Comments  []Comment `gorm:"foreignKey:AuthorID" json:"-"`
	Followers []Follow  `gorm:"foreignKey:FollowingID" json:"-"`
	Following []Follow  `gorm:"foreignKey:FollowerID" json:"-"`
}

// ─── Follow ───────────────────────────────────────────────────────────────────

type Follow struct {
	gorm.Model
	FollowerID  uint `gorm:"not null;index" json:"follower_id"`
	FollowingID uint `gorm:"not null;index" json:"following_id"`

	Follower  User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}

// ─── Tag ─────────────────────────────────────────────────────────────────────

type Tag struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Description string `gorm:"size:200" json:"description"`
	Color       string `gorm:"size:20;default:'#6366f1'" json:"color"`
	PostCount   int    `gorm:"default:0" json:"post_count"`

	Posts []Post `gorm:"many2many:post_tags" json:"-"`
}

// ─── Post ─────────────────────────────────────────────────────────────────────

type PostType string
type PostStatus string

const (
	PostTypePost    PostType = "post"
	PostTypeArticle PostType = "article"
	PostTypeTopic   PostType = "topic"
)

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusHidden    PostStatus = "hidden"
)

type Post struct {
	gorm.Model
	Title     string     `gorm:"not null;size:200" json:"title"`
	Content   string     `gorm:"not null;type:text" json:"content"`
	Summary   string     `gorm:"size:500" json:"summary"`
	Cover     string     `gorm:"size:500" json:"cover"`
	Type      PostType   `gorm:"type:varchar(20);default:'post'" json:"type"`
	Status    PostStatus `gorm:"type:varchar(20);default:'published'" json:"status"`
	AuthorID  uint       `gorm:"not null;index" json:"author_id"`
	ViewCount int        `gorm:"default:0" json:"view_count"`
	LikeCount int        `gorm:"default:0" json:"like_count"`
	PinTop    bool       `gorm:"default:false" json:"pin_top"`

	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags" json:"tags,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"-"`
	Likes    []Like    `gorm:"foreignKey:PostID" json:"-"`
}

// ─── Comment ──────────────────────────────────────────────────────────────────

type Comment struct {
	gorm.Model
	Content  string `gorm:"not null;type:text" json:"content"`
	PostID   uint   `gorm:"not null;index" json:"post_id"`
	AuthorID uint   `gorm:"not null;index" json:"author_id"`
	ParentID *uint  `gorm:"index" json:"parent_id"`
	LikeCount int   `gorm:"default:0" json:"like_count"`

	Post     Post      `gorm:"foreignKey:PostID" json:"-"`
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Parent   *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies  []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Likes    []Like    `gorm:"foreignKey:CommentID" json:"-"`
}

// ─── Like ─────────────────────────────────────────────────────────────────────

type Like struct {
	gorm.Model
	UserID    uint  `gorm:"not null;index" json:"user_id"`
	PostID    *uint `gorm:"index" json:"post_id"`
	CommentID *uint `gorm:"index" json:"comment_id"`

	User    User     `gorm:"foreignKey:UserID" json:"-"`
	Post    *Post    `gorm:"foreignKey:PostID" json:"-"`
	Comment *Comment `gorm:"foreignKey:CommentID" json:"-"`
}

// ─── Notification ─────────────────────────────────────────────────────────────

type NotificationType string

const (
	NotifyComment  NotificationType = "comment"
	NotifyLike     NotificationType = "like"
	NotifyFollow   NotificationType = "follow"
	NotifyReply    NotificationType = "reply"
	NotifySystem   NotificationType = "system"
)

type Notification struct {
	gorm.Model
	UserID     uint             `gorm:"not null;index" json:"user_id"`
	SenderID   *uint            `gorm:"index" json:"sender_id"`
	Type       NotificationType `gorm:"type:varchar(30)" json:"type"`
	Content    string           `gorm:"size:500" json:"content"`
	TargetID   *uint            `json:"target_id"`
	TargetType string           `gorm:"size:50" json:"target_type"`
	IsRead     bool             `gorm:"default:false" json:"is_read"`

	User   User  `gorm:"foreignKey:UserID" json:"-"`
	Sender *User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

// ─── SignIn (daily check-in) ──────────────────────────────────────────────────

type SignIn struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	SignDate  time.Time `gorm:"not null" json:"sign_date"`
	Score     int       `gorm:"default:5" json:"score"`
	Continued int       `gorm:"default:1" json:"continued"`
}

// ─── Report ───────────────────────────────────────────────────────────────────

type ReportStatus string

const (
	ReportPending  ReportStatus = "pending"
	ReportResolved ReportStatus = "resolved"
	ReportRejected ReportStatus = "rejected"
)

type Report struct {
	gorm.Model
	ReporterID  uint         `gorm:"not null;index" json:"reporter_id"`
	TargetID    uint         `gorm:"not null" json:"target_id"`
	TargetType  string       `gorm:"size:50;not null" json:"target_type"` // post | comment | user
	Reason      string       `gorm:"size:500;not null" json:"reason"`
	Status      ReportStatus `gorm:"type:varchar(20);default:'pending'" json:"status"`
	HandlerID   *uint        `json:"handler_id"`
	HandleNote  string       `gorm:"size:500" json:"handle_note"`

	Reporter User  `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
	Handler  *User `gorm:"foreignKey:HandlerID" json:"handler,omitempty"`
}
