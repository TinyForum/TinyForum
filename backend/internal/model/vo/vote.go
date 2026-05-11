package vo

import "time"

// VoteVO 投票记录脱敏视图（对外暴露）
type VoteVO struct {
	ID        uint      `json:"id"`
	CommentID uint      `json:"comment_id"`
	Value     int       `json:"value"` // 1: 赞同, -1: 反对, 0: 取消
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
