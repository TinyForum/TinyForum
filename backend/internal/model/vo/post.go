package vo

import "time"

// PostVO 帖子脱敏视图（对外暴露）
type PostVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title            string `json:"title"`
	Content          string `json:"content"`
	Summary          string `json:"summary,omitempty"`
	Cover            string `json:"cover,omitempty"`
	Type             string `json:"type"`              // PostType: post, question, etc.
	PostStatus       string `json:"post_status"`       // draft, published...
	ModerationStatus string `json:"moderation_status"` // normal, pending, rejected...

	AuthorID uint `json:"author_id"`
	// 脱敏后的作者信息
	Author struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar,omitempty"`
	} `json:"author,omitempty"`

	BoardID uint `json:"board_id"`
	// 可选：版块简要信息（需额外查询）
	Board struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"board,omitempty"`

	Tags []string `json:"tags,omitempty"` // 标签名称列表

	ViewCount  int  `json:"view_count"`
	LikeCount  int  `json:"like_count"`
	PinTop     bool `json:"pin_top"`
	PinInBoard bool `json:"pin_in_board"`
}
