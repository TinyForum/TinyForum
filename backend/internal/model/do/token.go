// model/token.go
package do

import (
	"time"
)

// RefreshToken 刷新令牌表
type RefreshToken struct {
	BaseModel
	UserID    uint      `gorm:"index:idx_user_id;not null" json:"user_id"`            // 用户ID
	Token     string    `gorm:"uniqueIndex:idx_token;size:500;not null" json:"token"` // 令牌
	JTI       string    `gorm:"uniqueIndex:idx_jti;size:36;not null" json:"jti"`      // JWT ID，用于精确撤销
	UserAgent string    `gorm:"size:500" json:"user_agent"`                           // 设备信息
	IP        string    `gorm:"size:45" json:"ip"`                                    // 登录IP
	ExpiresAt time.Time `gorm:"index:idx_expires_at" json:"expires_at"`               // 过期时间
	IsUsed    bool      `gorm:"default:false" json:"is_used"`                         // 是否已使用
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
