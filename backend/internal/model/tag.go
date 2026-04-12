package model

type Tag struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;not null;size:50" json:"name"`
	Description string `gorm:"size:200" json:"description"`
	Color       string `gorm:"size:20;default:'#6366f1'" json:"color"`
	PostCount   int    `gorm:"default:0" json:"post_count"`

	Posts []Post `gorm:"many2many:post_tags" json:"-"`
}
