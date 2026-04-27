package tag

// TODO: Refactory，不符合 REATful 规范
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
// base URL: /api/v1
func (h *TagHandler) RegisterRoutes(api *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	tagGroup := api.Group("/tags")
	{
		// 公开路由
		tagGroup.GET("", h.List)    // GET /api/v1/tags
		tagGroup.GET("/:id", h.Get) // GET /api/v1/tags/:id

		// 需要管理员权限的路由
		adminGroup := tagGroup.Group("")
		adminGroup.Use(mw.AuthMW(), mw.AdminRequiredMW())
		{
			adminGroup.POST("", h.Create)       // POST /api/v1/tags
			adminGroup.PUT("/:id", h.Update)    // PUT /api/v1/tags/:id
			adminGroup.DELETE("/:id", h.Delete) // DELETE /api/v1/tags/:id
		}
	}

}
