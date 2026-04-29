package answer

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *AnswerHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	answerGroup := api.Group("/answers")
	{
		answerGroup.GET("/:id", mw.OptionalAuth(), h.GetAnswer)
		answerGroup.DELETE("/:id", mw.Auth(), h.DeleteAnswer)
		answerGroup.GET("/:id/status", mw.OptionalAuth(), h.GetVoteStatus)
		answerGroup.POST("/:id/vote", mw.OptionalAuth(), h.VoteAnswer)
		answerGroup.DELETE("/:id/vote", mw.Auth(), h.RemoveVote)
		answerGroup.POST("/:id/accept", mw.Auth(), h.AcceptAnswer)
		answerGroup.POST("/:id/unaccept", mw.Auth(), h.UnacceptAnswer)
	}
}
