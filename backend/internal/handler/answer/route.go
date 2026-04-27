package answer

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *AnswerHandler) RegisterRoutes(api *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	answerGroup := api.Group("/answers")
	{
		answerGroup.GET("/:id", mw.OptionalAuthMW(), h.GetAnswer)
		answerGroup.DELETE("/:id", mw.AuthMW(), h.DeleteAnswer)
		answerGroup.GET("/:id/status", mw.OptionalAuthMW(), h.GetVoteStatus)
		answerGroup.POST("/:id/vote", mw.OptionalAuthMW(), h.VoteAnswer)
		answerGroup.DELETE("/:id/vote", mw.AuthMW(), h.RemoveVote)
		answerGroup.POST("/:id/accept", mw.AuthMW(), h.AcceptAnswer)
		answerGroup.POST("/:id/unaccept", mw.AuthMW(), h.UnacceptAnswer)
	}
}
