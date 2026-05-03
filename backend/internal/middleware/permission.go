// internal/middleware/permission.go
//
// 修复说明：
//   原版用 do.UserRole 断言 context 里的 user_role，但 Auth 中间件注入的是
//   string（来自 JWT claims.Role），导致断言永远失败，返回 500。
//   修复：先读 string，再转 do.UserRole，类型安全且向后兼容。

package middleware

import (
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// getRoleFromContext 从 gin context 安全读取 user_role，统一转为 do.UserRole。
// Auth 中间件注入的是 string（来自 JWT claims.Role），此处做一次转换。
func getRoleFromContext(c *gin.Context) (do.UserRole, bool) {
	raw, exists := c.Get(ContextUserRole)
	if !exists {
		return "", false
	}
	// Auth 注入的是 string
	if s, ok := raw.(string); ok {
		return do.UserRole(s), true
	}
	// 兼容直接注入 do.UserRole 的场景（单测 mock 等）
	if r, ok := raw.(do.UserRole); ok {
		return r, true
	}
	return "", false
}

// RequirePermission 检查用户是否拥有指定权限（基于角色的静态权限矩阵）。
// 适用于与具体资源无关的全局权限判断，如"是否能创建帖子"。
func RequirePermission(perm do.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := getRoleFromContext(c)
		if !exists {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}
		if !do.HasPermission(userRole, perm) {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAllPermissions 要求用户同时拥有全部指定权限。
func RequireAllPermissions(perms ...do.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := getRoleFromContext(c)
		if !exists {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}
		if !do.HasAllPermissions(userRole, perms...) {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAnyPermission 要求用户拥有任意一个指定权限。
func RequireAnyPermission(perms ...do.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := getRoleFromContext(c)
		if !exists {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}
		if !do.HasAnyPermission(userRole, perms...) {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireRole 要求用户角色优先级不低于 target。
// 例：RequireRole(do.RoleModerator) 允许 moderator / admin / super_admin 通过。
func RequireRole(target do.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := getRoleFromContext(c)
		if !exists {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}
		if !do.IsRoleAtLeast(userRole, target) {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}
