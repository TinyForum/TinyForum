package model

type PostType string
type PostStatus string

const (
	PostTypePost    PostType = "post"
	PostTypeArticle PostType = "article"
	PostTypeTopic   PostType = "topic"
)

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPending   PostStatus = "pending" // 待审核（命中 review 级敏感词 或 举报聚合触发）
	PostStatusPublished PostStatus = "published"
	PostStatusHidden    PostStatus = "hidden" // 已隐藏（审核拒绝 或 管理员操作）
)

type Post struct {
	BaseModel
	Title     string     `gorm:"not null;size:200" json:"title"`
	Content   string     `gorm:"not null;type:text" json:"content"`
	Summary   string     `gorm:"size:500" json:"summary"`
	Cover     string     `gorm:"size:500" json:"cover"`
	Type      PostType   `gorm:"type:varchar(20);default:'post'" json:"type"`
	Status    PostStatus `gorm:"type:varchar(20);default:'published'" json:"status"`
	AuthorID  uint       `gorm:"not null;index" json:"author_id"`
	ViewCount int        `gorm:"default:0" json:"view_count"`
	LikeCount int        `gorm:"default:0" json:"like_count"`
	PinTop    bool       `gorm:"default:false" json:"pin_top"`

	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags" json:"tags,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"-"`
	Likes    []Like    `gorm:"foreignKey:PostID" json:"-"`

	BoardID    uint      `gorm:"index" json:"board_id"`
	PinInBoard bool      `gorm:"default:false" json:"pin_in_board"`
	Board      Board     `gorm:"foreignKey:BoardID" json:"board,omitempty"`
	Question   *Question `gorm:"foreignKey:PostID" json:"question,omitempty"`
}
