package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

type SignIn struct {
	common.BaseModel
	UserID    uint      `gorm:"not null;index" json:"user_id"` // 用户id
	SignDate  time.Time `gorm:"not null" json:"sign_date"`     // 签到日期
	Score     int       `gorm:"default:5" json:"score"`        // 签到积分
	Continued int       `gorm:"default:1" json:"continued"`    // 连续签到天数
}
