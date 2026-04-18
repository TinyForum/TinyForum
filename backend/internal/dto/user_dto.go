// dto/leaderboard.go
package dto

// SimpleLeaderboardItem 精简版（仅核心字段）
type SimpleLeaderboardItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}

// DetailLeaderboardItem 详细版（包含头像等信息）
// type DetailLeaderboardItem struct {
// 	ID       uint   `json:"id"`
// 	Username string `json:"username"`
// 	Avatar   string `json:"avatar"`
// 	Score    int    `json:"score"`
// 	Rank     int    `json:"rank"`
// }

type LeaderboardUserSimple struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
}
type LeaderboardUserDetail struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Rank     int    `json:"rank"`
}

type LeaderboardRequest struct {
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
	// Fields string `form:"fields"`
}

// LeaderboardItemResponse 排行榜条目响应
type LeaderboardItemResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}

type LeaderboardResponse struct {
	Items []LeaderboardItemResponse `json:"items"`
}
