package topic

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *TopicHandler) RegisterRoutes(api *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	topicGroup := api.Group("/topics")
	{
		topicGroup.GET("", h.List)
		topicGroup.GET("/:id", h.GetByID)
		topicGroup.GET("/:id/posts", h.GetTopicPosts)
		topicGroup.GET("/:id/followers", h.GetFollowers)
		topicGroup.POST("", mw.AuthMW(), h.Create)
		topicGroup.PUT("/:id", mw.AuthMW(), h.Update)
		topicGroup.DELETE("/:id", mw.AuthMW(), h.Delete)
		topicGroup.POST("/:id/posts", mw.AuthMW(), h.AddPost)
		topicGroup.DELETE("/:id/posts/:post_id", mw.AuthMW(), h.RemovePost)
		topicGroup.POST("/:id/follow", mw.AuthMW(), h.Follow)
		topicGroup.DELETE("/:id/follow", mw.AuthMW(), h.Unfollow)
	}
}
