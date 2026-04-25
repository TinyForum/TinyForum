package middleware

import (
	"fmt"
	"tiny-forum/internal/model"
	riskservice "tiny-forum/internal/service/risk"
	"tiny-forum/pkg/ratelimit"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RateLimitMiddleware 返回一个基于用户风险等级的滑动窗口限流中间件
//
// 用法：
//
//	router.POST("/posts", RateLimitMiddleware(db, riskSvc, ratelimit.ActionCreatePost), handler)
func RateLimitMiddleware(db *gorm.DB, riskSvc riskservice.RiskService, action ratelimit.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDRaw, exists := c.Get(ContextUserID)
		if !exists {
			c.Next()
			return
		}

		userID, ok := userIDRaw.(uint)
		if !ok {
			c.Next()
			return
		}

		var user model.User
		if err := db.Select("id, score, is_blocked, created_at").First(&user, userID).Error; err != nil {
			c.Next()
			return
		}

		if user.IsBlocked {
			response.Forbidden(c, "账号已被封禁")
			c.Abort()
			return
		}

		result, err := riskSvc.CheckRateLimit(c.Request.Context(), &user, action)
		if err != nil {
			c.Next()
			return
		}

		if !result.Allowed {
			c.Header("X-RateLimit-Limit", fmt.Sprint(result.Limit))
			c.Header("X-RateLimit-Current", fmt.Sprint(result.Current))
			c.Header("X-RateLimit-Reset", fmt.Sprint(int(result.ResetIn.Seconds())))
			response.TooManyRequests(c, fmt.Sprintf("操作过于频繁，请 %d 分钟后再试", int(result.ResetIn.Minutes())+1))
			c.Abort()
			return
		}

		c.Next()
	}
}
