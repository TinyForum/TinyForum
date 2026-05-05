package do

import "tiny-forum/internal/model/common"

type CommentStatus string

const (
	CommentStatusVisible CommentStatus = "visible" // 正常可见
	CommentStatusPending CommentStatus = "pending" // 待审核
	CommentStatusHidden  CommentStatus = "hidden"  // 已隐藏
)

type Comment struct {
	common.BaseModel
	Content   string        `gorm:"not null;type:text" json:"content"`
	PostID    uint          `gorm:"not null;index" json:"post_id"`
	AuthorID  uint          `gorm:"not null;index" json:"author_id"`
	ParentID  *uint         `gorm:"index" json:"parent_id"`
	LikeCount int           `gorm:"default:0" json:"like_count"`
	Status    CommentStatus `gorm:"type:varchar(20);default:'visible';index" json:"status"` // 新增审核状态

	Post       Post      `gorm:"foreignKey:PostID" json:"-"`
	Author     User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Parent     *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Replies    []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
	Likes      []Like    `gorm:"foreignKey:CommentID" json:"-"`
	IsAnswer   bool      `gorm:"default:false;index" json:"is_answer"`
	IsAccepted bool      `gorm:"default:false;index" json:"is_accepted"`
	VoteCount  int       `gorm:"default:0" json:"vote_count"`
}
