package post

import "tiny-forum/internal/model/do"

type CreatePostInput struct {
	Title   string        `json:"title" binding:"required,min=2,max=200"`
	Content string        `json:"content" binding:"required,min=10"`
	Summary string        `json:"summary"`
	Cover   string        `json:"cover"`
	Type    string        `json:"type"`
	BoardID uint          `json:"board_id" binding:"required"`
	TagIDs  []uint        `json:"tag_ids"`
	Status  do.PostStatus `json:"status"`
}

type UpdatePostInput struct {
	Title   string `json:"title" binding:"min=2,max=200"`
	Content string `json:"content" binding:"min=10"`
	Summary string `json:"summary"`
	Cover   string `json:"cover"`
	TagIDs  []uint `json:"tag_ids"`
}
