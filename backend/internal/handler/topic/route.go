package topic

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *TopicHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	topicGroup := api.Group("/topics")
	{
		topicGroup.GET("", h.List)
		topicGroup.GET("/:id", h.GetByID)
		topicGroup.GET("/:id/posts", h.GetTopicPosts)
		topicGroup.GET("/:id/followers", h.GetFollowers)
		topicGroup.POST("", mw.Auth(), h.Create)
		topicGroup.PUT("/:id", mw.Auth(), h.Update)
		topicGroup.DELETE("/:id", mw.Auth(), h.Delete)
		topicGroup.POST("/:id/posts", mw.Auth(), h.AddPost)
		topicGroup.DELETE("/:id/posts/:post_id", mw.Auth(), h.RemovePost)
		topicGroup.POST("/:id/follow", mw.Auth(), h.Follow)
		topicGroup.DELETE("/:id/follow", mw.Auth(), h.Unfollow)
	}
}
