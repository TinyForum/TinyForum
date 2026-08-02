package do

// 话题下的文章结构
// type TopicPostsListItems struct {
// 	common.BaseModel
// 	TopicID   uint `gorm:"not null;uniqueIndex:idx_topic_post" json:"topic_id"`
// 	PostID    uint `gorm:"not null;uniqueIndex:idx_topic_post" json:"post_id"`
// 	SortOrder int  `gorm:"default:0" json:"sort_order"`
// 	AddedBy   uint `json:"added_by"`

// 	Topic   Topic   `gorm:"foreignKey:TopicID" json:"-"`
// 	Article Article `gorm:"foreignKey:CreationsID" json:"post,omitempty"`
// }

// // 表名
// func (TopicPostsListItems) TableName() string {
// 	return "topic_posts_list_items"
// }

import "tiny-forum/internal/model/common"

// Topic 观点，表达用户对某事的立场
// 主题下的作品列表项
type TopicCreation struct {
	common.BaseModel
	// Description string `gorm:"size:500" json:"description"`     // 描述
	TopicID    uint     `gorm:"primaryKey;index"`                          // 关联的话题 ID
	CreationID uint     `gorm:"uniqueIndex;not null"  json:"creations_id"` // 外键，唯一索引保证一对一
	Creation   Creation `gorm:"foreignKey:CreationID;references:ID" json:"creation,omitempty"`
	SeriesID   uint     `gorm:"index" json:"series_id"`   // 所属系列（可选）
	Viewpoint  string   `gorm:"size:50" json:"viewpoint"` // 观点倾向（支持/中立/反对）

	IsPublic      bool `gorm:"default:true;index" json:"is_public"` // 是否公开
	PostCount     int  `gorm:"default:0" json:"post_count"`         // 帖子数量
	FollowerCount int  `gorm:"default:0" json:"follower_count"`      // 关注者数量
	SortOrder     int  `gorm:"default:0"`                           // 排序顺序
	CreatorID     uint `gorm:"not null;index"`                      // 添加者ID

	// Creator User `gorm:"foreignKey:CreatorID" json:"creator,omitempty"` // 创建者
	// FollowersID []int `gorm:"foreignKey:TopicID" json:"-"` // 关注者ID
	// Creation *Creation `gorm:"foreignKey:CreationID" json:"-"` // 关联的Creation
}

func (TopicCreation) TableName() string {
	return "topic_creations"
}
