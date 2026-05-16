package do

import "tiny-forum/internal/model/common"

type Follow struct {
	common.BaseModel
	// 移除 uniqueIndex 标签，只保留普通索引
	FollowerID  uint `gorm:"not null;index:idx_follower;constraint:OnDelete:CASCADE"`
	FollowingID uint `gorm:"not null;index:idx_following;constraint:OnDelete:CASCADE"`

	Follower  *User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`
	Following *User `gorm:"foreignKey:FollowingID" json:"following,omitempty"`
}
