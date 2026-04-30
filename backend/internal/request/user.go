package request

type LeaderboardRequest struct {
	Limit int `form:"limit,default=20" binding:"min=1,max=100"`
	// Fields string `form:"fields"`
}
