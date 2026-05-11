package vo

import "time"

// FollowVO 关注关系脱敏视图（对外暴露）
type FollowVO struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	FollowerID  uint      `json:"follower_id"`
	FollowingID uint      `json:"following_id"`
}
