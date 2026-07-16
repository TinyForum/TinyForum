package question

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *QuestionHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	questionGroup := api.Group("/questions")
	{
		questionGroup.POST("/create", mw.Auth(), h.CreateQuestion)                                    // 创建问题
		questionGroup.GET("/simple", h.GetQuestionSimple)                                             // 获取问题列表
		questionGroup.GET("/list", mw.OptionalAuth(), h.GetQuestionsList)                             // 获取问题列表（分页）
		questionGroup.GET("/detail/:id", mw.OptionalAuth(), h.GetQuestionDetail)                      // 获取问题详情
		questionGroup.POST("/:post_id/answers", mw.Auth(), h.answerHandler.CreateAnswer)              // 创建回答
		questionGroup.GET("/:post_id/answers", mw.OptionalAuth(), h.answerHandler.GetQuestionAnswers) // 获取问题的回答列表
	}
}
