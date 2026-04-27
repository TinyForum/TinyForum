package question

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *QuestionHandler) RegisterRoutes(api *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	questionGroup := api.Group("/questions")
	{
		questionGroup.GET("/simple", h.GetQuestionSimple)
		questionGroup.GET("/list", mw.OptionalAuthMW(), h.GetQuestionsList)
		questionGroup.POST("/create", mw.AuthMW(), h.CreateQuestion)
		questionGroup.GET("/detail/:id", mw.OptionalAuthMW(), h.GetQuestionDetail)
		questionGroup.GET("/:id/answers", mw.OptionalAuthMW(), h.answerHandler.GetQuestionAnswers)
		questionGroup.POST("/:id/answers", mw.AuthMW(), h.answerHandler.CreateAnswer)
	}
}
