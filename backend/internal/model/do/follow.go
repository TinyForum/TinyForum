package do

import "tiny-forum/internal/model/common"

type Follow struct {
	common.BaseModel
	FollowerID  uint `gorm:"not null;index" json:"follower_id"`
	FollowingID uint `gorm:"not null;index" json:"following_id"`

	Follower  User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}
