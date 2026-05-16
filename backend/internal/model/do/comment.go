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
	Content   string        `gorm:"not null;type:text" json:"content"`                      // 评论内容
	PostID    uint          `gorm:"not null;index" json:"post_id"`                          // 关联的帖子ID
	AuthorID  uint          `gorm:"not null;index" json:"author_id"`                        // 评论作者ID
	ParentID  *uint         `gorm:"index" json:"parent_id"`                                 // 父评论ID
	LikeCount int           `gorm:"default:0" json:"like_count"`                            // 点赞数
	Status    CommentStatus `gorm:"type:varchar(20);default:'visible';index" json:"status"` // 审核状态

	Post       Post      `gorm:"foreignKey:PostID" json:"-"`                   // 关联的帖子
	Author     User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`  // 评论作者
	Parent     *Comment  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`  // 父评论
	Replies    []Comment `gorm:"foreignKey:ParentID" json:"replies,omitempty"` // 子评论
	Likes      []Like    `gorm:"-" json:"-"`                                   // 点赞记录
	IsAnswer   bool      `gorm:"default:false;index" json:"is_answer"`         // 是否为答案
	IsAccepted bool      `gorm:"default:false;index" json:"is_accepted"`       // 是否被采纳
	VoteCount  int       `gorm:"default:0" json:"vote_count"`                  // 投票数
}
