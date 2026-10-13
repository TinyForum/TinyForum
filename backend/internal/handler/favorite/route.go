package favorite

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册收藏相关路由
func (h *FavoriteHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	favs := api.Group("/users/me/favorites", mw.Auth())
	{
		favs.GET("", h.ListMyFavorites)              // GET    /api/v1/users/me/favorites       我的收藏
		favs.POST("", h.AddFavorite)                 // POST   /api/v1/users/me/favorites       收藏内容
		favs.DELETE("/:target_id", h.RemoveFavorite) // DELETE /api/v1/users/me/favorites/:target_id 取消收藏
	}
}
