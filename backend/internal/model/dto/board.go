package dto

import "time"

type GetBoardPostsRequest struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1,max=100"`
}

type GetBoardPostsResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	Cover      string    `json:"cover"`
	Type       string    `json:"type"`
	AuthorID   uint      `json:"author_id"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
}
