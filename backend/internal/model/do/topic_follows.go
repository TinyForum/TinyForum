package do

import "tiny-forum/internal/model/common"

type TopicFollow struct {
	common.BaseModel
	UserID  uint `gorm:"not null;uniqueIndex:idx_user_topic" json:"user_id"`
	TopicID uint `gorm:"not null;uniqueIndex:idx_user_topic" json:"topic_id"`

	User  User   `gorm:"foreignKey:UserID" json:"-"`
	Topic *Topic `gorm:"foreignKey:TopicID" json:"-"`
}

// 表名
func (TopicFollow) TableName() string {
	return "topic_follows"
}
