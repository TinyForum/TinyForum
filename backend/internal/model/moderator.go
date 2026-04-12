package model

import "encoding/json"

type Moderator struct {
	BaseModel
	UserID             uint            `gorm:"not null;uniqueIndex:idx_user_board" json:"user_id"`
	BoardID            uint            `gorm:"not null;uniqueIndex:idx_user_board" json:"board_id"`
	Permissions        json.RawMessage `gorm:"type:json" json:"permissions"`
	CanDeletePost      bool            `gorm:"default:false" json:"can_delete_post"`
	CanPinPost         bool            `gorm:"default:false" json:"can_pin_post"`
	CanEditAnyPost     bool            `gorm:"default:false" json:"can_edit_any_post"`
	CanManageModerator bool            `gorm:"default:false" json:"can_manage_moderator"`
	CanBanUser         bool            `gorm:"default:false" json:"can_ban_user"`

	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Board Board `gorm:"foreignKey:BoardID" json:"board,omitempty"`
}
