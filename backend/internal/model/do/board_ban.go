package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

type BoardBan struct {
	common.BaseModel
	// 唯一索引必须是部分索引（deleted_at IS NULL），否则软删后可解封记录会阻塞再次封禁。
	// 与 migrations/20261013_board_ban_partial_unique.up.sql 保持一致。
	UserID    uint       `gorm:"not null;uniqueIndex:idx_active_ban,priority:1,where:deleted_at IS NULL" json:"user_id"`  // 被封禁用户
	BoardID   uint       `gorm:"not null;uniqueIndex:idx_active_ban,priority:2" json:"board_id"` // 被禁言板块
	BannedBy  uint       `gorm:"not null" json:"banned_by"`                                      // 封禁者
	Reason    string     `gorm:"type:text" json:"reason"`                                        // 禁言原因
	ExpiresAt *time.Time `gorm:"index" json:"expires_at"`                                        // 封禁到期时间

	User   *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`     // 被禁言用户
	Board  *Board `gorm:"foreignKey:BoardID" json:"board,omitempty"`   // 被禁言板块
	Banner *User  `gorm:"foreignKey:BannedBy" json:"banner,omitempty"` // 封禁者
}
