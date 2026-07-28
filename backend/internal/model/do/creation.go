package do

import (
	"tiny-forum/internal/model/common"

	"gorm.io/datatypes"
)

// Creation 内容基表，所有类型公共字段
type Creation struct {
	common.BaseModel
	Title   string `gorm:"not null;size:200" json:"title"`    // 标题
	Content string `gorm:"not null;type:text" json:"content"` // 内容
	Summary string `gorm:"size:500" json:"summary"`           // 摘要
	Slug    string `gorm:"size:180;uniqueIndex" json:"slug"`  // URL标识

	CoverUrl  string                      `gorm:"size:500" json:"cover_url"`                            // 封面
	ImageUrls datatypes.JSONSlice[string] `gorm:"type:text" json:"image_urls"`                          // 图片
	Type      CreationType                `gorm:"type:varchar(20);not null;index" json:"creation_type"` // 类型

	// 作品状态（用户主动状态）
	CreationStatus CreationStatus `gorm:"type:varchar(20);default:'draft';index" json:"creation_status"` // 状态
	// 审核状态（被动）
	ModerationStatus ModerationStatus `gorm:"type:varchar(20);default:'normal';index" json:"moderation_status"` // 审核状态

	AuthorID   uint `gorm:"not null;index" json:"author_id"` // 作者
	ViewCount  int  `gorm:"default:0" json:"view_count"`     // 浏览数
	LikeCount  int  `gorm:"default:0" json:"like_count"`     // 点赞数
	PinTop     bool `gorm:"default:false" json:"pin_top"`    // 置顶
	IsOriginal bool `gorm:"default:true" json:"is_original"` // 是否原创

	BoardID    uint `gorm:"index" json:"board_id"`             // 所属板块
	PinInBoard bool `gorm:"default:false" json:"pin_in_board"` // 板块置顶

	// 关联（子表一对一）
	Article       *Article       `gorm:"foreignKey:ID;references:ID" json:"article,omitempty"`
	Question      *Question      `gorm:"foreignKey:ID;references:ID" json:"question,omitempty"`
	Post          *Post          `gorm:"foreignKey:ID;references:ID" json:"post,omitempty"`
	TopicCreation *TopicCreation `gorm:"foreignKey:ID;references:ID" json:"topic_creation,omitempty"`

	// 其他关联（多对多、一对多）
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`   // 作者
	Board    Board     `gorm:"foreignKey:BoardID" json:"board,omitempty"`     // 所属板块
	Tags     []Tag     `gorm:"many2many:creation_tags" json:"tags,omitempty"` // 标签
	Comments []Comment `gorm:"foreignKey:WorksID;references:ID" json:"-"`
	Likes    []Like    `gorm:"-" json:"-"` // 点赞
}

func (Creation) TableName() string {
	return "creations"
}
