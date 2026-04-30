package vo

// LeaderboardItemResponse 排行榜条目响应
//
//	type LeaderboardItemVO struct {
//		ID       uint   `json:"id"`
//		Username string `json:"username"`
//		Avatar   string `json:"avatar"`
//		Score    int    `json:"score"`
//		Rank     int    `json:"rank"`
//	}
//
// SimpleLeaderboardItem 精简版（仅核心字段）
type SimpleLeaderboardItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}

type Statistics struct {
	TotalPosts          int64 `json:"total_posts"`
	TotalComments       int64 `json:"total_comments"`
	TotalFavorites      int64 `json:"total_favorites"`
	UnreadNotifications int64 `json:"unread_notifications"`
}
