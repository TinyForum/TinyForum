package do

import "tiny-forum/internal/model/common"

type Tag struct {
	common.BaseModel
	Name        string `gorm:"uniqueIndex;not null;size:50" json:"name"` // 标签名称
	Description string `gorm:"size:200" json:"description"`              // 标签描述
	Color       string `gorm:"size:20;default:'#6366f1'" json:"color"`   // 标签颜色
	PostCount   int    `gorm:"default:0" json:"post_count"`              // 标签关联的文章数量

	Posts []Post `gorm:"many2many:post_tags" json:"-"` // 标签关联的文章
}

// 表名
func (Tag) TableName() string {
	return "tags"
}
