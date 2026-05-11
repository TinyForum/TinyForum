package vo

import "time"

// LikeVO 点赞记录脱敏视图（对外暴露）
type LikeVO struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id"`
	PostID    *uint     `json:"post_id,omitempty"`
	CommentID *uint     `json:"comment_id,omitempty"`
}
