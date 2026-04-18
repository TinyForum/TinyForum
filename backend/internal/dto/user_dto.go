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
type DetailLeaderboardItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}
