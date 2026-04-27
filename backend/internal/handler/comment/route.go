package comment

// TODO: Refactory，不符合 REATful 规范
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
// base URL: /api/v1
func (h *CommentHandler) RegisterRoutes(api *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	commentGroup := api.Group("/comments")
	{
		commentGroup.GET("/post/:post_id", h.List)
		commentGroup.POST("",
			mw.AuthMW(),
			mw.RateLimitMW("create_comment"),
			mw.ContentCheckMW([]string{"content"}),
			h.Create,
		)
		commentGroup.DELETE("/:id", mw.AuthMW(), h.Delete)
	}

}
