package answer

import (
	articleService "tiny-forum/internal/service/article"
	commentService "tiny-forum/internal/service/comment"
	questionService "tiny-forum/internal/service/question"
)

// AnswerHandler 处理回答相关请求
type AnswerHandler struct {
	questionSvc questionService.QuestionService
	commentSvc  commentService.CommentService
	postSvc     articleService.ArticleService
}

// NewAnswerHandler 创建 AnswerHandler 实例
func NewAnswerHandler(
	questionSvc questionService.QuestionService,
	commentSvc commentService.CommentService,
	postSvc articleService.ArticleService,
) *AnswerHandler {
	return &AnswerHandler{
		questionSvc: questionSvc,
		commentSvc:  commentSvc,
		postSvc:     postSvc,
	}
}
