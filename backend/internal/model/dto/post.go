package dto

import (
	"time"
	"tiny-forum/internal/model/po"
)

// PostListOptions 帖子列表查询选项
type PostListOptions struct {
	AuthorID         uint
	TagID            uint
	PostType         string
	Keyword          string
	SortBy           string
	Status           po.PostStatus
	ModerationStatus po.ModerationStatus
}

// GetStatsDay 获取每日统计数据
type StatsDayQuery struct {
	Date string `form:"date" binding:"omitempty,datetime=2006-01-02"`
	Type string `form:"type" binding:"omitempty,oneof=users posts comments all"`
}

// GetStatsResponse 响应的统计数据
type GetStatsResponse struct {
	Day time.Time `json:"day"`
}
