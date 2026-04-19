package dto

import "tiny-forum/internal/model"

// PostListOptions 帖子列表查询选项
type PostListOptions struct {
	AuthorID         uint
	TagID            uint
	PostType         string
	Keyword          string
	SortBy           string                 // "" = latest, "hot" = popular
	Status           model.PostStatus       // "" = all, "pending" = pending, "approved" = approved, "rejected" = rejected
	ModerationStatus model.ModerationStatus // "" = all, "low" = low, "medium" = medium, "high" = high
}
