package vo

import "time"

// UserFavorite 我的收藏视图对象
type UserFavorite struct {
	ID        uint      `json:"id"`         // 收藏记录ID
	CreatedAt time.Time `json:"created_at"` // 收藏时间
	TargetID  uint      `json:"target_id"`  // 作品（帖子）ID
	Title     string    `json:"title"`      // 标题
	Summary   string    `json:"summary"`    // 摘要
	LikeCount int       `json:"like_count"` // 点赞数
	BoardName string    `json:"board_name"` // 所属板块
}
