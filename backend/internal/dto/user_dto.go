// dto/leaderboard.go
package dto

// DetailLeaderboardItem 详细版（包含头像等信息）
// type DetailLeaderboardItem struct {
// 	ID       uint   `json:"id"`
// 	Username string `json:"username"`
// 	Avatar   string `json:"avatar"`
// 	Score    int    `json:"score"`
// 	Rank     int    `json:"rank"`
// }

type LeaderboardUserDetail struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Rank     int    `json:"rank"`
}

type LeaderboardUserSimple struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
}

// type LeaderboardResponse struct {
// 	Items []LeaderboardItemResponse `json:"items"`
// }

type GlobalStatsCount struct {
	TotalPosts     int `json:"total_posts"`
	TotalComments  int `json:"total_comments"`
	TotalFavorites int `json:"total_favorites"`
	TotalViolation int `json:"total_violation"`
}
