package comment

import (
	commentService "tiny-forum/internal/service/comment"
	questionService "tiny-forum/internal/service/question"
)

type CommentHandler struct {
	commentSvc  commentService.CommentService
	questionSvc *questionService.QuestionService
}

func NewCommentHandler(commentSvc commentService.CommentService, questionSvc *questionService.QuestionService) *CommentHandler {
	return &CommentHandler{
		commentSvc:  commentSvc,
		questionSvc: questionSvc,
	}
}
