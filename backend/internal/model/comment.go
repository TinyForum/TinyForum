package model

type Comment struct {
	BaseModel
	Content   string `gorm:"not null;type:text" json:"content"`
	PostID    uint   `gorm:"not null;index" json:"post_id"`
	AuthorID  uint   `gorm:"not null;index" json:"author_id"`
	ParentID  *uint  `gorm:"index" json:"parent_id"`
	LikeCount int    `gorm:"default:0" json:"like_count"`

	Post       Post      `gorm:"foreignKey:PostID" json:"-"`
	Author     User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Parent     *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies    []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Likes      []Like    `gorm:"foreignKey:CommentID" json:"-"`
	IsAnswer   bool      `gorm:"default:false;index" json:"is_answer"`   // 新增
	IsAccepted bool      `gorm:"default:false;index" json:"is_accepted"` // 新增
	VoteCount  int       `gorm:"default:0" json:"vote_count"`            // 新增
}
