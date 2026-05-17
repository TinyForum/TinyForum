package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

// BlockedIP IP封禁记录
type BlockedIP struct {
	common.BaseModel
	IP         string     `gorm:"type:varchar(45);uniqueIndex;not null"` // 支持IPv6
	Reason     string     `gorm:"type:text"`                             // 封禁原因
	OperatorID uint       `gorm:"index"`                                 // 操作员ID
	ExpireAt   *time.Time `gorm:"index"`                                 // 封禁过期时间，NULL表示永久封禁
}

func (BlockedIP) TableName() string {
	return "blocked_ips"
}

// IsExpired 检查封禁是否已过期
func (b *BlockedIP) IsExpired() bool {
	if b.ExpireAt == nil {
		return false
	}
	return time.Now().After(*b.ExpireAt)
}

// IsPermanent 是否是永久封禁
func (b *BlockedIP) IsPermanent() bool {
	return b.ExpireAt == nil
}
