package question

import (
	"tiny-forum/internal/handler/answer"
	commentService "tiny-forum/internal/service/comment"
	postService "tiny-forum/internal/service/post"
	questionService "tiny-forum/internal/service/question"
)

type QuestionHandler struct {
	questionSvc   questionService.QuestionService
	commentSvc    commentService.CommentService
	postSvc       postService.PostService
	answerHandler *answer.AnswerHandler
}

func NewQuestionHandler(
	questionSvc questionService.QuestionService,
	commentSvc commentService.CommentService,
	postSvc postService.PostService,
	answerHandler *answer.AnswerHandler,
) *QuestionHandler {
	return &QuestionHandler{
		questionSvc:   questionSvc,
		commentSvc:    commentSvc,
		postSvc:       postSvc,
		answerHandler: answerHandler, // ✅ 修正拼写
	}
}
