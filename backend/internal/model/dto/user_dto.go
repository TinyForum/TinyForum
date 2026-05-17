// dto/leaderboard.go
package dto

type LeaderboardUserDetail struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	Score     int    `json:"score"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Rank      int    `json:"rank"`
}

type LeaderboardUserSimple struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	AvatarUrl string `json:"avatar_url"`
	Score     int    `json:"score"`
}

// type LeaderboardResponse struct {
// 	Items []LeaderboardItemResponse `json:"items"`
// }

type GlobalStatsCount struct {
	TotalCountPosts     int `json:"total_count_posts"`
	TotalCountComments  int `json:"total_count_comments"`
	TotalCountFavorites int `json:"total_count_favorites"`
	TotalCountViolation int `json:"total_count_violation"`
}

type UserStatsCount struct {
	TotalCountPosts     int `json:"total_count_posts"`
	TotalCountComments  int `json:"total_count_comments"`
	TotalCountFavorites int `json:"total_count_favorites"`
	TotalCountViolation int `json:"total_count_violation"`
}
