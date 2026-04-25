package question

import (
	commentService "tiny-forum/internal/service/comment"
	postService "tiny-forum/internal/service/post"
	questionService "tiny-forum/internal/service/question"
)

type QuestionHandler struct {
	questionSvc questionService.QuestionService
	commentSvc  commentService.CommentService
	postSvc     postService.PostService
}

func NewQuestionHandler(
	questionSvc questionService.QuestionService,
	commentSvc commentService.CommentService,
	postSvc postService.PostService,
) *QuestionHandler {
	return &QuestionHandler{
		questionSvc: questionSvc,
		commentSvc:  commentSvc,
		postSvc:     postSvc,
	}
}
