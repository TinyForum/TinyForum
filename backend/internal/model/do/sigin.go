package do

import (
	"time"
	"tiny-forum/internal/model/common"
)

type SignIn struct {
	common.BaseModel
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	SignDate  time.Time `gorm:"not null" json:"sign_date"`
	Score     int       `gorm:"default:5" json:"score"`
	Continued int       `gorm:"default:1" json:"continued"`
}
