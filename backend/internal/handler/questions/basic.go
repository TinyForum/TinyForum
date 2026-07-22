package question

import (
	"tiny-forum/internal/handler/answer"
	articleService "tiny-forum/internal/service/article"
	commentService "tiny-forum/internal/service/comment"
	questionService "tiny-forum/internal/service/question"
)

type QuestionHandler struct {
	questionSvc   questionService.QuestionService
	commentSvc    commentService.CommentService
	postSvc       articleService.ArticleService
	answerHandler *answer.AnswerHandler
}

func NewQuestionHandler(
	questionSvc questionService.QuestionService,
	commentSvc commentService.CommentService,
	postSvc articleService.ArticleService,
	answerHandler *answer.AnswerHandler,
) *QuestionHandler {
	return &QuestionHandler{
		questionSvc:   questionSvc,
		commentSvc:    commentSvc,
		postSvc:       postSvc,
		answerHandler: answerHandler, // ✅ 修正拼写
	}
}
