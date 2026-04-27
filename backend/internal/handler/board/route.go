package board

// TODO: Refactory，不符合 REATful 规范
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func (h *BoardHandler) RegisterRoutes(board *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	// 用户排行榜
	g := board.Group("/slug")
	{
		g.GET("/:slug", h.GetBoardBySlug)                            // 获取版块信息
		g.GET("/:slug/posts", mw.OptionalAuthMW(), h.GetPostsBySlug) // 获取版块帖子
	}
	g = board.Group("/:id")
	{
		g.GET("", mw.OptionalAuthMW(), h.GetByID)                            // 获取版块信息
		g.POST("/moderators/apply-moderator", mw.AuthMW(), h.ApplyModerator) // 申请成为版主
		g.PUT("", mw.AuthMW(), mw.AdminRequiredMW(), h.Update)               // 更新版块信息
		g.DELETE("/:id", mw.AuthMW(), mw.AdminRequiredMW(), h.Delete)        // 删除版块

	}

	g = board.Group("/moderators")
	{
		g.GET("/apply-status", mw.AuthMW(), h.GetUserApplications)
		g.GET("/managed", mw.AuthMW(), h.GetUserModeratorBoards)
	}
}
