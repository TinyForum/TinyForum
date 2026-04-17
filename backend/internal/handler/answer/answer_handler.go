package answer

import (
	commentService "tiny-forum/internal/service/comment"
	postService "tiny-forum/internal/service/post"
	questionService "tiny-forum/internal/service/question"
)

// AnswerHandler 处理回答相关请求
type AnswerHandler struct {
	questionSvc *questionService.QuestionService
	commentSvc  *commentService.CommentService
	postSvc     *postService.PostService
}

// NewAnswerHandler 创建 AnswerHandler 实例
func NewAnswerHandler(
	questionSvc *questionService.QuestionService,
	commentSvc *commentService.CommentService,
	postSvc *postService.PostService,
) *AnswerHandler {
	return &AnswerHandler{
		questionSvc: questionSvc,
		commentSvc:  commentSvc,
		postSvc:     postSvc,
	}
}
