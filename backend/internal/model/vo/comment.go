package vo

import "time"

// CommentVO 评论脱敏视图（对外暴露）
type CommentVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Content   string    `json:"content"`

	AuthorID uint `json:"author_id"`
	// 可选：脱敏后的作者信息（需额外查询填充）
	Author struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar,omitempty"`
	} `json:"author,omitempty"`

	PostID   uint  `json:"post_id"`
	ParentID *uint `json:"parent_id,omitempty"`

	LikeCount  int    `json:"like_count"`
	VoteCount  int    `json:"vote_count"`
	Status     string `json:"status,omitempty"` // CommentStatus 映射为字符串
	IsAnswer   bool   `json:"is_answer"`
	IsAccepted bool   `json:"is_accepted"`
}
