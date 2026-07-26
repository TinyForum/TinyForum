package bo

import "tiny-forum/internal/model/do"

// PostListOptions 帖子列表查询选项
type ListPosts struct {
	AuthorID         uint
	TagNames         []string
	Type             string
	Keyword          string
	SortBy           string
	PostStatus       do.CreationStatus
	ModerationStatus do.ModerationStatus
}

// page, pageSize int, keyword string
