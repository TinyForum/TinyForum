package do

import (
	"fmt"
	"time"
	"tiny-forum/internal/model/common"
)

// UserTimeline 用户各时间线的最后阅读时间
type UserTimeline struct {
	common.BaseModel
	UserID       uint         `gorm:"not null;uniqueIndex:idx_user_timeline,priority:1" json:"user_id"`               // 用户ID
	TimelineType TimelineType `gorm:"type:varchar(20);uniqueIndex:idx_user_timeline,priority:2" json:"timeline_type"` // 时间线类型
	LastReadAt   time.Time    `json:"last_read_at"`                                                                   // 最后阅读时间
}

func (UserTimeline) TableName() string {
	return "user_timelines"
}

// TimelineType 时间线类型
type TimelineType string

const (
	TimelineHome   TimelineType = "home"   // 首页时间线
	TimelineFollow TimelineType = "follow" // 关注动态
	TimelineMine   TimelineType = "mine"   // 个人动态
)

// enum [home, follow, mine]

func (t TimelineType) IsValid() bool {
	switch t {
	case TimelineHome, TimelineFollow, TimelineMine:
		return true
	}
	return false
}

func ParseTimelineType(s string) (TimelineType, error) {
	tt := TimelineType(s)
	if tt.IsValid() {
		return tt, nil
	}
	return "", fmt.Errorf("invalid timeline type: %s", s)
}
