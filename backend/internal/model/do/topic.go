package do

import "tiny-forum/internal/model/common"

type Topic struct {
	common.BaseModel
	Title         string `gorm:"not null;size:150;uniqueIndex;comment:话题标题" json:"title"` // 话题标题
	Slug          string `gorm:"size:180;uniqueIndex;comment:URL标识" json:"slug"`          // URL标识
	Description   string `gorm:"size:500;comment:话题描述" json:"description"`                // 话题描述
	CoverUrl      string `gorm:"size:500;comment:封面图URL" json:"cover_url"`                // 封面图URL
	CreatorID     uint   `gorm:"not null;index;comment:创建者ID" json:"creator_id"`          // 创建者ID
	IsPublic      bool   `gorm:"default:true;index;comment:是否公开" json:"is_public"`        // 是否公开
	PostCount     int    `gorm:"default:0;comment:帖子数量" json:"post_count"`                // 帖子数量
	FollowerCount int    `gorm:"default:0;comment:关注者数量" json:"follower_count"`           // 关注者数量

	Creator   User          `gorm:"foreignKey:CreatorID" json:"creator,omitempty"` // 创建者
	Posts     []TopicPost   `gorm:"foreignKey:TopicID" json:"-"`                   // 帖子
	Followers []TopicFollow `gorm:"foreignKey:TopicID" json:"-"`                   // 关注者
}

// 表名
// TableName 方法用于定义数据库表名
func (Topic) TableName() string {
	// 返回 topics 作为数据库表名
	return "topics"
}
