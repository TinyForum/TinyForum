package vo

type SimpleLeaderboardItem struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Score    int    `json:"score"`
	Rank     int    `json:"rank"`
}
