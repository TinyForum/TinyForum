package middleware

import (
	"strings"

	"tiny-forum/pkg/jwt"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "user_id"
	ContextUsername = "username"
	ContextUserRole = "user_role"
)

func Auth(jwtMgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "无效的认证格式")
			c.Abort()
			return
		}

		claims, err := jwtMgr.Parse(parts[1])
		if err != nil {
			response.Unauthorized(c, "Token 无效或已过期，请重新登录")
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
		if role != "admin" {
			response.Forbidden(c, "需要管理员权限")
			c.Abort()
			return
		}
		c.Next()
	}
}

// OptionalAuth sets user info if token is present, but doesn't block if not
func OptionalAuth(jwtMgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				if claims, err := jwtMgr.Parse(parts[1]); err == nil {
					c.Set(ContextUserID, claims.UserID)
					c.Set(ContextUsername, claims.Username)
					c.Set(ContextUserRole, claims.Role)
				}
			}
		}
		c.Next()
	}
}
