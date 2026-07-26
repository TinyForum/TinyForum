package do

import (
	"tiny-forum/internal/model/common"

	"gorm.io/datatypes"
)

// Post 短文（微博风格），无额外字段，仅关联 Creation
type Post struct {
	common.BaseModel
	CreationID    uint                        `gorm:"uniqueIndex;not null"  json:"creation_id"` // 关联 Creation
	AllowComments bool                        `gorm:"default:true" json:"allow_comments"`       // 是否允许评论
	ImageUrls     datatypes.JSONSlice[string] `gorm:"type:text" json:"image_urls"`              // 图片链接列表
	Creation      *Creation                   `gorm:"foreignKey:CreationID" json:"-"`           // 关联 Creation
}

func (Post) TableName() string {
	return "posts"
}
