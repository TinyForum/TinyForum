package do

import "tiny-forum/internal/model/common"

type Like struct {
	common.BaseModel
	UserID    uint  `gorm:"not null;index" json:"user_id"`
	PostID    *uint `gorm:"index" json:"post_id"`
	CommentID *uint `gorm:"index" json:"comment_id"`

	User    User     `gorm:"foreignKey:UserID" json:"-"`
	Post    *Post    `gorm:"foreignKey:PostID" json:"-"`
	Comment *Comment `gorm:"foreignKey:CommentID" json:"-"`
}
