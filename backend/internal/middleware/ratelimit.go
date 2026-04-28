package middleware

import (
	"fmt"
	"log"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/model"
	riskservice "tiny-forum/internal/service/risk"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RateLimitMiddleware 返回一个基于用户风险等级的滑动窗口限流中间件
func RateLimitMiddleware(db *gorm.DB, riskSvc riskservice.RiskService, action ratelimit.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 添加调试日志
		log.Printf("[RateLimit] Processing request: %s %s, action: %s",
			c.Request.Method, c.Request.URL.Path, action)

		userIDRaw, exists := c.Get(ContextUserID)

		// 情况1：未登录用户，使用IP限流
		if !exists {
			handleAnonymousUser(c, riskSvc, action)
			return
		}

		// 情况2：已登录用户，使用用户ID限流
		handleAuthenticatedUser(c, db, riskSvc, action, userIDRaw)
	}
}

// handleAnonymousUser 处理未登录用户的限流
func handleAnonymousUser(c *gin.Context, riskSvc riskservice.RiskService, action ratelimit.Action) {
	ip := c.ClientIP()
	log.Printf("[RateLimit] Anonymous user, IP: %s, action: %s", ip, action)

	// IP白名单检查
	if isWhitelistedIP(ip) {
		log.Printf("[RateLimit] IP %s is whitelisted, skipping rate limit", ip)
		c.Next()
		return
	}

	result, err := riskSvc.CheckRateLimitByIP(c.Request.Context(), ip, action)
	if err != nil {
		log.Printf("[RateLimit] CheckRateLimitByIP error for IP %s: %v, allowing request", ip, err)
		c.Next()
		return
	}

	log.Printf("[RateLimit] IP %s - Allowed: %v, Current: %d, Limit: %d, ResetIn: %.0fs",
		ip, result.Allowed, result.Current, result.Limit, result.ResetIn)

	resetSeconds := int(result.ResetIn.Seconds())
	if resetSeconds < 0 {
		resetSeconds = 0
	}
	if !result.Allowed {
		setRateLimitHeaders(c, result)
		response.TooManyRequests(c, fmt.Sprintf("操作过于频繁，请 %d 秒后再试", resetSeconds))
		c.Abort()
		return
	}

	// 添加限流信息到响应头
	if result.Limit > 0 {
		c.Header("X-RateLimit-Remaining", fmt.Sprint(result.Limit-result.Current))
	}

	c.Next()
}

// handleAuthenticatedUser 处理已登录用户的限流
func handleAuthenticatedUser(c *gin.Context, db *gorm.DB, riskSvc riskservice.RiskService,
	action ratelimit.Action, userIDRaw interface{}) {

	userID, ok := userIDRaw.(uint)
	if !ok {
		log.Printf("[RateLimit] Invalid user ID type, skipping rate limit")
		c.Next()
		return
	}

	log.Printf("[RateLimit] Authenticated user, UserID: %d, action: %s", userID, action)

	var user model.User
	if err := db.Select("id, score, is_blocked, created_at").First(&user, userID).Error; err != nil {
		log.Printf("[RateLimit] User not found: %d, error: %v, fallback to IP rate limit", userID, err)
		// 用户不存在，降级为IP限流
		handleAnonymousUser(c, riskSvc, action)
		return
	}

	if user.IsBlocked {
		log.Printf("[RateLimit] User %d is blocked", userID)
		response.Forbidden(c, "账号已被封禁")
		c.Abort()
		return
	}

	result, err := riskSvc.CheckRateLimit(c.Request.Context(), &user, action)
	if err != nil {
		log.Printf("[RateLimit] CheckRateLimit error for user %d: %v, fallback to IP rate limit", userID, err)
		// 降级为IP限流
		handleAnonymousUser(c, riskSvc, action)
		return
	}

	log.Printf("[RateLimit] User %d - Allowed: %v, Current: %d, Limit: %d, ResetIn: %.0fs",
		userID, result.Allowed, result.Current, result.Limit, result.ResetIn)

	resetSeconds := int(result.ResetIn.Seconds())
	if resetSeconds < 0 {
		resetSeconds = 0
	}
	if !result.Allowed {
		setRateLimitHeaders(c, result)

		// 根据用户信用分显示不同的提示
		if user.Score < 50 {
			response.TooManyRequests(c, "您的账号信用分较低，操作频率受限，请提升信用分后再试")
		} else {
			response.TooManyRequests(c, fmt.Sprintf("操作过于频繁，请 %d 秒后再试", resetSeconds))
		}
		c.Abort()
		return
	}

	c.Next()
}

// setRateLimitHeaders 设置限流相关的响应头
func setRateLimitHeaders(c *gin.Context, result ratelimit.Result) {
	c.Header("X-RateLimit-Limit", fmt.Sprint(result.Limit))
	c.Header("X-RateLimit-Current", fmt.Sprint(result.Current))
	c.Header("X-RateLimit-Reset", fmt.Sprint(int(result.ResetIn)))
	c.Header("Retry-After", fmt.Sprint(int(result.ResetIn)))
}

// isWhitelistedIP IP白名单检查
func isWhitelistedIP(ip string) bool {
	whitelist := []string{
		// "127.0.0.1",
		// "::1",
		// 添加其他白名单IP
	}
	for _, allowed := range whitelist {
		if ip == allowed {
			return true
		}
	}
	return false
}
