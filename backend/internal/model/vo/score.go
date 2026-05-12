package vo

// UserScoreResponse 积分响应（用于列表）
type UserScoreVO struct {
	ID uint `json:"id"`
	// Username string `json:"username"`
	// Avatar   string `json:"avatar_url"`
	Score int `json:"score"`
}
