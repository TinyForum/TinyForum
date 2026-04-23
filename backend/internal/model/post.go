package model

type PostType string

const (
	PostTypePost     PostType = "post"
	PostTypeArticle  PostType = "article"
	PostTypeTopic    PostType = "topic"
	PostTypeQuestion PostType = "question"
)

// 合法的帖子类型集合
var validPostTypes = map[PostType]bool{
	PostTypePost:    true,
	PostTypeArticle: true,
	PostTypeTopic:   true,
}

type PostStatus string

// const (
// 	PostTypePost    PostType = "post"
// 	PostTypeArticle PostType = "article"
// 	PostTypeTopic   PostType = "topic"
// )

// 用户主动控制的状态（用户能感知、能操作）
const (
	// 草稿（用户保存未发布）
	PostStatusDraft PostStatus = "draft"
	// 待用户确认/提交（如编辑后重新提交）
	PostStatusPending PostStatus = "pending"
	// 已发布（用户主动发布）
	PostStatusPublished PostStatus = "published"
	// 用户隐藏（如自己删除/隐藏，或管理员操作但以用户视角展示）
	PostStatusHidden PostStatus = "hidden"
)

// 系统风控状态（由内容安全模块自动判定或管理员审核结果）

type Post struct {
	BaseModel
	Title   string   `gorm:"not null;size:200" json:"title"`
	Content string   `gorm:"not null;type:text" json:"content"`
	Summary string   `gorm:"size:500" json:"summary"`
	Cover   string   `gorm:"size:500" json:"cover"`
	Type    PostType `gorm:"type:varchar(20);default:'post'" json:"type"`

	// 用户状态 - 用户自己控制/感知的状态
	PostStatus PostStatus `gorm:"type:varchar(20);default:'draft'" json:"post_status"`

	// 系统风控状态 - 由风控引擎或管理员审核决定
	ModerationStatus ModerationStatus `gorm:"type:varchar(20);default:'normal'" json:"moderation_status"`

	AuthorID  uint `gorm:"not null;index" json:"author_id"`
	ViewCount int  `gorm:"default:0" json:"view_count"`
	LikeCount int  `gorm:"default:0" json:"like_count"`
	PinTop    bool `gorm:"default:false" json:"pin_top"`

	// 关联
	Author   User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Tags     []Tag     `gorm:"many2many:post_tags" json:"tags,omitempty"`
	Comments []Comment `gorm:"foreignKey:PostID" json:"-"`
	Likes    []Like    `gorm:"foreignKey:PostID" json:"-"`

	BoardID    uint      `gorm:"index" json:"board_id"`
	PinInBoard bool      `gorm:"default:false" json:"pin_in_board"`
	Board      Board     `gorm:"foreignKey:BoardID" json:"board,omitempty"`
	Question   *Question `gorm:"foreignKey:PostID" json:"question,omitempty"`
}

// IsValid 检查帖子类型是否合法
func (pt PostType) IsValid() bool {
	return validPostTypes[pt]
}

// 可选：从字符串安全转换
func ParsePostType(s string) PostType {
	pt := PostType(s)
	if pt.IsValid() {
		return pt
	}
	return PostTypePost // 默认值
}
