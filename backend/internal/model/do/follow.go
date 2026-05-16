package do

import "tiny-forum/internal/model/common"

type Follow struct {
	common.BaseModel
	FollowerID  uint `gorm:"not null;index:idx_follower;uniqueIndex:idx_follow_unique,priority:1;constraint:OnDelete:CASCADE"`  // 关注者
	FollowingID uint `gorm:"not null;index:idx_following;uniqueIndex:idx_follow_unique,priority:2;constraint:OnDelete:CASCADE"` // 被关注者

	// 使用指针类型，避免 JSON 空对象
	Follower  *User `gorm:"foreignKey:FollowerID" json:"follower,omitempty"`   // 关注者
	Following *User `gorm:"foreignKey:FollowingID" json:"following,omitempty"` // 被关注者
}
