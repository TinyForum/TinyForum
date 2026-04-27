package post

// TODO: Refactory，不符合 REATful 规范
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
// base URL: /api/v1
func (h *PostHandler) RegisterRoutes(api *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	postGroup := api.Group("/posts")
	{
		postGroup.GET("", mw.OptionalAuthMW(), h.List) // 用户获取帖子列表
		postGroup.GET("/:id", mw.OptionalAuthMW(), h.GetByID)
		postGroup.POST("",
			mw.AuthMW(),
			mw.RateLimitMW("create_post"),
			mw.ContentCheckMW([]string{"title", "content"}),
			h.Create,
		) // 用户发布帖子
		postGroup.PUT("/:id", mw.AuthMW(), h.Update)
		postGroup.DELETE("/:id", mw.AuthMW(), h.Delete)
		postGroup.POST("/:id/like", mw.AuthMW(), h.Like)
		postGroup.DELETE("/:id/like", mw.AuthMW(), h.Unlike)
	}

}
