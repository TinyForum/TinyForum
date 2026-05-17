package do

import "tiny-forum/internal/model/common"

// TimelineSubscription 用户订阅关系
type TimelineSubscription struct {
	common.BaseModel
	SubscriberID uint          `gorm:"not null;uniqueIndex:idx_subscription,priority:1" json:"subscriber_id"`       // 订阅者ID
	TargetType   SubTargetType `gorm:"type:varchar(20);uniqueIndex:idx_subscription,priority:2" json:"target_type"` // 目标类型
	TargetID     uint          `gorm:"not null;uniqueIndex:idx_subscription,priority:3" json:"target_id"`           // 目标ID（用户ID或板块ID）
	IsActive     bool          `gorm:"default:true" json:"is_active"`
}

// 表名
func (TimelineSubscription) TableName() string {
	return "timeline_subscriptions"
}

// SubTargetType 订阅目标类型
type SubTargetType string

const (
	SubTargetUser  SubTargetType = "user"  // 订阅用户
	SubTargetBoard SubTargetType = "board" // 订阅板块
)

// enum [user, board]

func (t SubTargetType) IsValid() bool {
	switch t {
	case SubTargetUser, SubTargetBoard:
		return true
	}
	return false
}

// 辅助方法：获取订阅的目标用户ID（当 TargetType == user 时）
func (s *TimelineSubscription) GetTargetUserID() uint {
	if s.TargetType == SubTargetUser {
		return s.TargetID
	}
	return 0
}
