package timeline

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *TimelineHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	timelineGroup := api.Group("/timeline", mw.Auth())
	{
		timelineGroup.GET("", h.GetHomeTimeline)                   // 获取首页时间线
		timelineGroup.GET("/following", h.GetFollowingTimeline)    // 获取关注用户的时间线
		timelineGroup.POST("/subscribe/:user_id", h.Subscribe)     // 订阅用户
		timelineGroup.DELETE("/subscribe/:user_id", h.Unsubscribe) // 取消订阅用户
		timelineGroup.GET("/subscriptions", h.GetSubscriptions)    // 获取订阅列表
	}
}
