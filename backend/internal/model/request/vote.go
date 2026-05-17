package request

import "tiny-forum/internal/model/do"

type VoteAnswerRequest struct {
	CommentID uint               `json:"comment_id" binding:"required"`
	VoteType  *do.AnswerVoteType `json:"vote_type" binding:"required,oneof=up down"`
}
