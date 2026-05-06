package request

type ListPosts struct {
	Page             int      `form:"page"`
	PageSize         int      `form:"page_size"`
	Keyword          string   `form:"keyword"`
	SortBy           string   `form:"sort_by"`
	PostType         string   `form:"type"`
	AuthorID         uint     `form:"author_id"`
	postTags         []string `form:"tags"`
	TagNames         []string `form:"tag_names"`
	PostStatus       string   `form:"status"`
	ModerationStatus string   `form:"moderation_status"`
}
