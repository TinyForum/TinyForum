package dto

import "tiny-forum/internal/model"

type CreateQuestionRequest struct {
	Title       string           `json:"title" binding:"required,max=100"`
	Content     string           `json:"content" binding:"required"`
	Summary     string           `json:"summary" binding:"max=500"`
	Cover       string           `json:"cover" binding:"omitempty,url"`
	BoardID     uint             `json:"board_id"`
	TagIDs      []uint           `json:"tag_ids"`
	RewardScore int              `json:"reward_score" binding:"min=0,max=100"`
	Status      model.PostStatus `json:"status"`
}
