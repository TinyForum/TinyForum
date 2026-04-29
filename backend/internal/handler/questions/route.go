package question

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *QuestionHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	questionGroup := api.Group("/questions")
	{
		questionGroup.GET("/simple", h.GetQuestionSimple)
		questionGroup.GET("/list", mw.OptionalAuth(), h.GetQuestionsList)
		questionGroup.POST("/create", mw.Auth(), h.CreateQuestion)
		questionGroup.GET("/detail/:id", mw.OptionalAuth(), h.GetQuestionDetail)
		questionGroup.GET("/:id/answers", mw.OptionalAuth(), h.answerHandler.GetQuestionAnswers)
		questionGroup.POST("/:id/answers", mw.Auth(), h.answerHandler.CreateAnswer)
	}
}
