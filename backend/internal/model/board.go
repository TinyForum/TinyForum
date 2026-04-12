package model

import "time"

type Board struct {
	BaseModel
	Name        string   `gorm:"not null;size:50;uniqueIndex" json:"name"`
	Slug        string   `gorm:"not null;size:50;uniqueIndex" json:"slug"`
	Description string   `gorm:"size:500" json:"description"`
	Icon        string   `gorm:"size:100" json:"icon"`
	Cover       string   `gorm:"size:500" json:"cover"`
	ParentID    *uint    `gorm:"index" json:"parent_id"`
	SortOrder   int      `gorm:"default:0" json:"sort_order"`
	ViewRole    UserRole `gorm:"type:varchar(20);default:'user'" json:"view_role"`
	PostRole    UserRole `gorm:"type:varchar(20);default:'user'" json:"post_role"`
	ReplyRole   UserRole `gorm:"type:varchar(20);default:'user'" json:"reply_role"`
	PostCount   int      `gorm:"default:0" json:"post_count"`
	ThreadCount int      `gorm:"default:0" json:"thread_count"`
	TodayCount  int      `gorm:"default:0" json:"today_count"`

	Parent     *Board      `gorm:"foreignKey:ParentID" json:"parent,omitempty" swaggerignore:"true"`
	Children   []Board     `gorm:"foreignKey:ParentID" json:"children,omitempty" swaggerignore:"true"`
	Moderators []Moderator `gorm:"foreignKey:BoardID" json:"-"`
}

// BoardTree 用于树形结构的响应（避免递归）
type BoardTree struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	Slug        string      `json:"slug"`
	Description string      `json:"description"`
	Icon        string      `json:"icon"`
	Cover       string      `json:"cover"`
	ParentID    *uint       `json:"parent_id"`
	SortOrder   int         `json:"sort_order"`
	ViewRole    UserRole    `json:"view_role"`
	PostRole    UserRole    `json:"post_role"`
	ReplyRole   UserRole    `json:"reply_role"`
	PostCount   int         `json:"post_count"`
	ThreadCount int         `json:"thread_count"`
	TodayCount  int         `json:"today_count"`
	Children    []BoardTree `json:"children,omitempty"`
}

type Moderator struct {
	BaseModel
	UserID             uint   `gorm:"not null;uniqueIndex:idx_user_board" json:"user_id"`
	BoardID            uint   `gorm:"not null;uniqueIndex:idx_user_board" json:"board_id"`
	Permissions        string `gorm:"type:json" json:"permissions"`
	CanDeletePost      bool   `gorm:"default:false" json:"can_delete_post"`
	CanPinPost         bool   `gorm:"default:false" json:"can_pin_post"`
	CanEditAnyPost     bool   `gorm:"default:false" json:"can_edit_any_post"`
	CanManageModerator bool   `gorm:"default:false" json:"can_manage_moderator"`
	CanBanUser         bool   `gorm:"default:false" json:"can_ban_user"`

	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Board Board `gorm:"foreignKey:BoardID" json:"board,omitempty"`
}

type BoardBan struct {
	BaseModel
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	BoardID   uint       `gorm:"not null;index" json:"board_id"`
	BannedBy  uint       `json:"banned_by"`
	Reason    string     `gorm:"size:500" json:"reason"`
	ExpiresAt *time.Time `json:"expires_at"`

	User   User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Board  Board `gorm:"foreignKey:BoardID" json:"board,omitempty"`
	Banner User  `gorm:"foreignKey:BannedBy" json:"banner,omitempty"`
}

type ModeratorLog struct {
	BaseModel
	ModeratorID uint   `gorm:"not null;index" json:"moderator_id"`
	BoardID     uint   `gorm:"index" json:"board_id"`
	Action      string `gorm:"type:varchar(50)" json:"action"`
	TargetType  string `gorm:"size:50" json:"target_type"`
	TargetID    uint   `json:"target_id"`
	Reason      string `gorm:"size:500" json:"reason"`
	OldValue    string `gorm:"type:json" json:"old_value"`
	NewValue    string `gorm:"type:json" json:"new_value"`

	Moderator User  `gorm:"foreignKey:ModeratorID" json:"moderator,omitempty"`
	Board     Board `gorm:"foreignKey:BoardID" json:"board,omitempty"`
}
