// model/token.go
package model

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken 刷新令牌表
type RefreshToken struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"index:idx_user_id;not null" json:"user_id"`
	Token     string         `gorm:"uniqueIndex:idx_token;size:500;not null" json:"token"`
	JTI       string         `gorm:"uniqueIndex:idx_jti;size:36;not null" json:"jti"` // JWT ID，用于精确撤销
	UserAgent string         `gorm:"size:500" json:"user_agent"`                      // 设备信息
	IP        string         `gorm:"size:45" json:"ip"`                               // 登录IP
	ExpiresAt time.Time      `gorm:"index:idx_expires_at" json:"expires_at"`          // 过期时间
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
