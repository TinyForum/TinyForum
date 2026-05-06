package dto

import (
	"time"
)

// GetStatsDay 获取每日统计数据
type StatsDayQuery struct {
	Date string `form:"date" binding:"omitempty,datetime=2006-01-02"`            // 日期
	Type string `form:"type" binding:"omitempty,oneof=users posts comments all"` // 类型
}

// GetStatsResponse 响应的统计数据
type GetStatsResponse struct {
	Day time.Time `json:"day"` // 日期
}

type Listposts struct {
}
