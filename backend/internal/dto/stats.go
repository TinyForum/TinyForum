package dto

// StatsRangeQuery 范围统计请求
type StatsRangeQuery struct {
	StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Type      string `form:"type" binding:"omitempty,oneof=users posts comments all"`
}

// DailyStat 每日统计数据
type DailyStat struct {
	Date       string `json:"date"` // 2026-04-21
	NewUser    int64  `json:"new_user"`
	NewArticle int64  `json:"new_article"`
	NewComment int64  `json:"new_comment"`
	NewBoard   int64  `json:"new_board"`
	NewTag     int64  `json:"new_tag"`
	ActiveUser int64  `json:"active_user"`
}
