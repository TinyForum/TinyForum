package timeline

import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
func (h *TimelineHandler) RegisterRoutes(api *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	timelineGroup := api.Group("/timeline", mw.AuthMW())
	{
		timelineGroup.GET("", h.GetHomeTimeline)
		timelineGroup.GET("/following", h.GetFollowingTimeline)
		timelineGroup.POST("/subscribe/:user_id", h.Subscribe)
		timelineGroup.DELETE("/subscribe/:user_id", h.Unsubscribe)
		timelineGroup.GET("/subscriptions", h.GetSubscriptions)
	}
}
