package answer

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *AnswerHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	answerGroup := api.Group("/answers")
	{
		answerGroup.GET("/:id", mw.OptionalAuth(), h.GetAnswer)            // 获取单个回答
		answerGroup.DELETE("/:id", mw.Auth(), h.DeleteAnswer)              // 删除单个回答
		answerGroup.GET("/:id/status", mw.OptionalAuth(), h.GetVoteStatus) // 获取单个回答的投票状态
		answerGroup.POST("/:id/vote", mw.OptionalAuth(), h.VoteAnswer)     // 投票单个回答
		answerGroup.DELETE("/:id/vote", mw.Auth(), h.RemoveVote)           // 取消投票单个回答
		answerGroup.POST("/:id/accept", mw.Auth(), h.AcceptAnswer)         // 接受单个回答
		answerGroup.POST("/:id/unaccept", mw.Auth(), h.UnacceptAnswer)     // 取消接受单个回答
	}
}
