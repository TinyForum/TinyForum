package tag

// TODO: Refactory，不符合 REATful 规范
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
// base URL: /api/v1
func (h *TagHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	tagGroup := api.Group("/tags")
	{
		// 公开路由
		tagGroup.GET("", h.List)    // GET /api/v1/tags 获取标签列表
		tagGroup.GET("/:id", h.Get) // GET /api/v1/tags/:id 获取单个标签

		// 需要管理员权限的路由
		adminGroup := tagGroup.Group("")
		adminGroup.Use(mw.Auth(), mw.AdminRequired())
		{
			adminGroup.POST("", h.Create)       // POST /api/v1/tags 创建标签
			adminGroup.PUT("/:id", h.Update)    // PUT /api/v1/tags/:id 更新标签
			adminGroup.DELETE("/:id", h.Delete) // DELETE /api/v1/tags/:id 删除标签
		}
	}

}
