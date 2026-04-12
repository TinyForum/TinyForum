package middleware

import (
	"tiny-forum/internal/model"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// RequirePermission 基于权限控制
func RequirePermission(perm model.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleInterface, exists := c.Get(ContextUserRole)
		if !exists {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		userRole, ok := userRoleInterface.(model.UserRole)
		if !ok {
			response.InternalError(c, "角色类型错误")
			c.Abort()
			return
		}

		if !model.HasPermission(userRole, perm) {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAllPermissions 需要多个权限（全部满足）
func RequireAllPermissions(perms ...model.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleInterface, exists := c.Get(ContextUserRole)
		if !exists {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		userRole, ok := userRoleInterface.(model.UserRole)
		if !ok {
			response.InternalError(c, "角色类型错误")
			c.Abort()
			return
		}

		if !model.HasAllPermissions(userRole, perms...) {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission 需要任一权限
func RequireAnyPermission(perms ...model.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleInterface, exists := c.Get(ContextUserRole)
		if !exists {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		userRole, ok := userRoleInterface.(model.UserRole)
		if !ok {
			response.InternalError(c, "角色类型错误")
			c.Abort()
			return
		}

		if !model.HasAnyPermission(userRole, perms...) {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}
