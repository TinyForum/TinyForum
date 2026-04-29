package post

// TODO: Refactory，不符合 REATful 规范
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
// base URL: /api/v1
func (h *PostHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	postGroup := api.Group("/posts")
	{
		postGroup.GET("", mw.OptionalAuth(), h.List) // 用户获取帖子列表
		postGroup.GET("/:id", mw.OptionalAuth(), h.GetByID)
		postGroup.POST("",
			mw.Auth(),
			mw.RateLimit("create_post"),
			mw.ContentCheck([]string{"title", "content"}),
			h.Create,
		) // 用户发布帖子
		postGroup.PUT("/:id", mw.Auth(), h.Update)
		postGroup.DELETE("/:id", mw.Auth(), h.Delete)
		postGroup.POST("/:id/like", mw.Auth(), h.Like)
		postGroup.DELETE("/:id/like", mw.Auth(), h.Unlike)
	}

}
