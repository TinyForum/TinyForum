package vo

import (
	"time"
)

// PostVO 帖子脱敏视图（对外暴露）
type ArticleVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CreationID uint       `gorm:"uniqueIndex;not null" json:"creations_id"` // 外键，唯一索引保证一对一
	Creation   CreationVO `gorm:"foreignKey:CreationID;references:ID" json:"creation,omitempty"`
	// 如果有 Article 特有字段，加在这里；若无，可保留空结构体
}
type PostListVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Summary   string    `json:"summary,omitempty"`
	Cover     string    `json:"cover,omitempty"`
	Author    struct {
		ID        uint   `json:"id"`
		Username  string `json:"username"`
		AvatarUrl string `json:"avatar_url,omitempty"`
	} `json:"author"`
}

type ArticleCardVO struct {
	Post  ArticleVO `json:"post"`
	Liked bool      `json:"liked"`
}
type BoardWithArticle struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
type AutherWithArticle struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url,omitempty"`
}

// HotArticleRow 热门文章查询结果行
type HotArticleRowVO struct {
	ID           int64
	Title        string
	BoardID      int64
	BoardName    string
	AuthorID     int64
	AuthorName   string
	ViewCount    int64
	CommentCount int64
	LikeCount    int64
}
