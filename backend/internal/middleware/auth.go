package middleware

import (
	"strings"
	"tiny-forum/internal/repository/token"
	"tiny-forum/pkg/jwt"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "user_id"
	ContextUsername = "username"
	ContextUserRole = "user_role"
)

// 从请求中提取 token，优先 Cookie，降级 Authorization Header
func extractToken(c *gin.Context) string {
	// 1. 优先读 HttpOnly Cookie
	if cookie, err := c.Cookie("tiny_forum_token"); err == nil && cookie != "" {
		return cookie
	}
	// 2. 降级读 Authorization Header（兼容直接调用 API 的场景）
	authHeader := c.GetHeader("Authorization")
	if parts := strings.SplitN(authHeader, " ", 2); len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}
	return ""
}

func Auth(jwtMgr *jwt.JWTManager, tokenRepo token.TokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c) // 从请求中提取 token
		if token == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}
		claims, err := jwtMgr.Parse(token)
		if err != nil {
			response.Unauthorized(c, "Token 无效或已过期，请重新登录")
			c.Abort()
			return
		}
		revoked, err := tokenRepo.IsTokenRevoked(c.Request.Context(), claims.ID)
		if err == nil && revoked {
			response.Unauthorized(c, "token 已失效，请重新登录")
			c.Abort()
			return
		}
		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextUserRole, claims.Role)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get(ContextUserRole)
		if role != "admin" && role != "super_admin" {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

func OptionalAuth(jwtMgr *jwt.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := extractToken(c); token != "" {
			if claims, err := jwtMgr.Parse(token); err == nil {
				c.Set(ContextUserID, claims.UserID)
				c.Set(ContextUsername, claims.Username)
				c.Set(ContextUserRole, claims.Role)
			}
		}
		c.Next()
	}
}
