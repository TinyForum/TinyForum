// internal/middleware/ratelimit.go

package middleware

import (
	"fmt"
	"log"
	"net"
	"strings"
	"tiny-forum/config"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/model/do"
	riskservice "tiny-forum/internal/service/risk"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RateLimitMiddleware 限流中间件结构体
type RateLimitMiddleware struct {
	db      *gorm.DB
	riskSvc riskservice.RiskService
	cfg     *config.RateLimitConfig
}

// NewRateLimitMiddleware 创建限流中间件实例
func NewRateLimitMiddleware(db *gorm.DB, riskSvc riskservice.RiskService, cfg *config.RateLimitConfig) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		db:      db,
		riskSvc: riskSvc,
		cfg:     cfg,
	}
}

// Middleware 返回限流中间件处理函数
func (m *RateLimitMiddleware) Middleware(action ratelimit.Action) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果限流未启用，直接放行
		if m.cfg != nil && !m.cfg.Enabled {
			c.Next()
			return
		}

		log.Printf("[RateLimit] Processing request: %s %s, action: %s",
			c.Request.Method, c.Request.URL.Path, action)

		userIDRaw, exists := c.Get(ContextUserID)

		// 情况1：未登录用户，使用IP限流
		if !exists {
			m.handleAnonymousUser(c, action)
			return
		}

		// 情况2：已登录用户，使用用户ID限流
		m.handleAuthenticatedUser(c, action, userIDRaw)
	}
}

// handleAnonymousUser 处理未登录用户的限流
func (m *RateLimitMiddleware) handleAnonymousUser(c *gin.Context, action ratelimit.Action) {
	ip := c.ClientIP()
	log.Printf("[RateLimit] Anonymous user, IP: %s, action: %s", ip, action)

	// 使用配置的白名单
	if m.isWhitelistedIP(ip) {
		log.Printf("[RateLimit] IP %s is whitelisted, skipping rate limit", ip)
		c.Next()
		return
	}

	result, err := m.riskSvc.CheckRateLimitByIP(c.Request.Context(), ip, action)
	if err != nil {
		log.Printf("[RateLimit] CheckRateLimitByIP error for IP %s: %v, allowing request", ip, err)
		c.Next()
		return
	}

	resetSeconds := max(int(result.ResetIn.Seconds()), 0)

	log.Printf("[RateLimit] IP %s - Allowed: %v, Current: %d, Limit: %d, ResetIn: %ds",
		ip, result.Allowed, result.Current, result.Limit, resetSeconds)

	if !result.Allowed {
		m.setRateLimitHeaders(c, result)
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
func (m *RateLimitMiddleware) handleAuthenticatedUser(c *gin.Context, action ratelimit.Action, userIDRaw any) {
	userID, ok := userIDRaw.(uint)
	if !ok {
		log.Printf("[RateLimit] Invalid user ID type, skipping rate limit")
		c.Next()
		return
	}

	log.Printf("[RateLimit] Authenticated user, UserID: %d, action: %s", userID, action)

	var user do.User
	if err := m.db.Select("id, score, is_blocked, created_at").First(&user, userID).Error; err != nil {
		log.Printf("[RateLimit] User not found: %d, error: %v, fallback to IP rate limit", userID, err)
		m.handleAnonymousUser(c, action)
		return
	}

	if user.IsBlocked {
		log.Printf("[RateLimit] User %d is blocked", userID)
		response.Forbidden(c, "账号已被封禁")
		c.Abort()
		return
	}

	result, err := m.riskSvc.CheckRateLimit(c.Request.Context(), &user, action)
	if err != nil {
		log.Printf("[RateLimit] CheckRateLimit error for user %d: %v, fallback to IP rate limit", userID, err)
		m.handleAnonymousUser(c, action)
		return
	}

	log.Printf("[RateLimit] User %d - Allowed: %v, Current: %d, Limit: %d, ResetIn: %ds",
		userID, result.Allowed, result.Current, result.Limit, result.ResetIn)

	resetSeconds := max(int(result.ResetIn.Seconds()), 0)
	if !result.Allowed {
		m.setRateLimitHeaders(c, result)

		if user.Score < 50 {
			response.TooManyRequests(c, "您的账号积分较低，操作频率受限，请提升积分后再试")
		} else {
			response.TooManyRequests(c, fmt.Sprintf("操作过于频繁，请 %d 秒后再试", resetSeconds))
		}
		c.Abort()
		return
	}

	c.Next()
}

// setRateLimitHeaders 设置限流相关的响应头
func (m *RateLimitMiddleware) setRateLimitHeaders(c *gin.Context, result ratelimit.Result) {
	c.Header("X-RateLimit-Limit", fmt.Sprint(result.Limit))
	c.Header("X-RateLimit-Current", fmt.Sprint(result.Current))
	c.Header("X-RateLimit-Reset", fmt.Sprint(int(result.ResetIn)))
	c.Header("Retry-After", fmt.Sprint(int(result.ResetIn)))
}

// isWhitelistedIP 检查 IP 是否在白名单中
func (m *RateLimitMiddleware) isWhitelistedIP(ip string) bool {
	if m.cfg == nil || len(m.cfg.IPWhitelist) == 0 {
		return false
	}

	for _, whitelistIP := range m.cfg.IPWhitelist {
		if strings.Contains(whitelistIP, "/") {
			// 支持 CIDR 格式
			if isIPInCIDR(ip, whitelistIP) {
				return true
			}
		} else if ip == whitelistIP {
			return true
		}
	}
	return false
}

// isIPInCIDR 检查 IP 是否在 CIDR 范围内
func isIPInCIDR(ip string, cidr string) bool {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return ipnet.Contains(parsedIP)
}
