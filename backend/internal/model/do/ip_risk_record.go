package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

// IPRiskRecord IP风险记录
type IPRiskRecord struct {
	common.BaseModel
	IP          string    `gorm:"type:varchar(45);index;not null"` // IP地址
	EventType   string    `gorm:"type:varchar(50);not null"`       // 事件类型
	EventDetail string    `gorm:"type:text"`                       // 事件详情
	ExpireAt    time.Time `gorm:"index"`                           // 过期时间
}

func (IPRiskRecord) TableName() string {
	return "ip_risk_records"
}
