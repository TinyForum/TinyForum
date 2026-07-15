package topic

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *TopicHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	topicGroup := api.Group("/topics")
	{
		topicGroup.GET("", h.List)                                        // 获取所有专题
		topicGroup.GET("/:id", h.GetByID)                                 // 获取单个专题
		topicGroup.GET("/:id/posts", h.GetTopicPosts)                     // 获取专题下的帖子
		topicGroup.GET("/:id/followers", h.GetFollowers)                  // 获取专题的关注者
		topicGroup.POST("", mw.Auth(), h.Create)                          // 创建专题
		topicGroup.PUT("/:id", mw.Auth(), h.Update)                       // 更新专题
		topicGroup.DELETE("/:id", mw.Auth(), h.Delete)                    // 删除专题
		topicGroup.POST("/:id/posts", mw.Auth(), h.AddPost)               // 向专题添加帖子
		topicGroup.DELETE("/:id/posts/:post_id", mw.Auth(), h.RemovePost) // 从专题移除帖子
		topicGroup.POST("/:id/follow", mw.Auth(), h.Follow)               // 关注专题
		topicGroup.DELETE("/:id/follow", mw.Auth(), h.Unfollow)           // 取消关注专题
	}
}
