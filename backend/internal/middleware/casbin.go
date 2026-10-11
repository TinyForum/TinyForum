// internal/middleware/casbin.go
//
// CasbinAuth 是路由级权限中间件，职责：
//   - 从 gin.Context 读取已由 Auth 中间件注入的 user_role（string）
//   - 未登录时 fallback 为 "guest"
//   - 调用 enforcer.Enforce(role, path, method) 决策
//   - 拒绝时返回 403，不泄露策略细节
//
// 使用方式（在路由注册时叠加在 Auth 之后）：
//
//	protected := api.Group("/posts", mw.Auth(), mw.CasbinAuth())
//	public     := api.Group("/posts", mw.OptionalAuth(), mw.CasbinAuth())  // guest 可读

package middleware

import (
	"tiny-forum/pkg/response"

	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

// casbinAuth 是真正执行 Enforce 的内部函数，供 middlewareSet 调用。
func casbinAuth(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := "guest"
		if raw, exists := c.Get(ContextUserRole); exists {
			if s, ok := raw.(string); ok && s != "" {
				role = s
			}
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		if path == "" {
			c.Next()
			return
		}

		ok, err := enforcer.Enforce(role, path, method)

		if err != nil {
			response.InternalError(c, "权限校验异常")
			c.Abort()
			return
		}
		if !ok {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}
