package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

type BoardBan struct {
	common.BaseModel
	UserID    uint       `gorm:"not null;uniqueIndex:idx_active_ban,priority:1" json:"user_id"`  // 被封禁用户
	BoardID   uint       `gorm:"not null;uniqueIndex:idx_active_ban,priority:2" json:"board_id"` // 被禁言板块
	BannedBy  uint       `gorm:"not null" json:"banned_by"`                                      // 封禁者
	Reason    string     `gorm:"type:text" json:"reason"`                                        // 禁言原因
	ExpiresAt *time.Time `gorm:"index" json:"expires_at"`                                        // 封禁到期时间

	User   *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`     // 被禁言用户
	Board  *Board `gorm:"foreignKey:BoardID" json:"board,omitempty"`   // 被禁言板块
	Banner *User  `gorm:"foreignKey:BannedBy" json:"banner,omitempty"` // 封禁者
}
