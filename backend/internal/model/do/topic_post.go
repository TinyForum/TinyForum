package do

import "tiny-forum/internal/model/common"

type TopicPost struct {
	common.BaseModel
	TopicID   uint `gorm:"not null;uniqueIndex:idx_topic_post" json:"topic_id"`
	PostID    uint `gorm:"not null;uniqueIndex:idx_topic_post" json:"post_id"`
	SortOrder int  `gorm:"default:0" json:"sort_order"`
	AddedBy   uint `json:"added_by"`

	Topic   Topic   `gorm:"foreignKey:TopicID" json:"-"`
	Article Article `gorm:"foreignKey:PostID" json:"post,omitempty"`
}

// 表名
func (TopicPost) TableName() string {
	return "topic_posts"
}
