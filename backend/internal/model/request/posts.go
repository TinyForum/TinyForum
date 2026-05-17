package request

import "tiny-forum/internal/model/do"

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

type CreatePostRequest struct {
	Title   string        `json:"title" binding:"required,min=2,max=200"`
	Content string        `json:"content" binding:"required,min=10"`
	Summary string        `json:"summary"`
	Cover   string        `json:"cover"`
	Type    string        `json:"type"`
	BoardID uint          `json:"board_id" binding:"required"`
	TagIDs  []uint        `json:"tag_ids"`
	Status  do.PostStatus `json:"status"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"min=2,max=200"`
	Content string `json:"content" binding:"min=10"`
	Summary string `json:"summary"`
	Cover   string `json:"cover"`
	TagIDs  []uint `json:"tag_ids"`
}
