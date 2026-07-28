package vo

type TopicCreationVO struct {
	Description string     `gorm:"size:500" json:"description"`               // 描述
	TopicID     uint       `gorm:"primaryKey;index"`                          // 关联的话题 ID
	CreationID  uint       `gorm:"uniqueIndex;not null"  json:"creations_id"` // 外键，唯一索引保证一对一
	Creation    CreationVO `gorm:"foreignKey:CreationID;references:ID" json:"creation,omitempty"`
	SeriesID    uint       `gorm:"index" json:"series_id"`   // 所属系列（可选）
	Viewpoint   string     `gorm:"size:50" json:"viewpoint"` // 观点倾向（支持/中立/反对）

	IsPublic      bool `gorm:"default:true;index" json:"is_public"` // 是否公开
	PostCount     int  `gorm:"default:0" json:"post_count"`         // 帖子数量
	FollowerCount int  `gorm:"defult:0" json:"follower_count"`      // 关注者数量
	SortOrder     int  `gorm:"default:0"`                           // 排序顺序
	CreatorID     uint `gorm:"not null;index"`                      // 添加者ID

}
