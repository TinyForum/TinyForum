package auth

// TODO: Refactory，不符合 REATful 规范
import (
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// 所有账号安全相关的都放在这里
// FIXME: UserHandler 中的账号安全功能将逐步迁移到 AuthHandler 中

// RegisterRoutes 注册路由
// base URL: /api/v1
// Group URL: /auth
func (h *AuthHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	// 用户排行榜
	g := api.Group("/auth")
	{
		g.POST("/register", mw.RateLimit(ratelimit.ActionRegister), h.Register) // 用户注册
		g.POST("/login", mw.RateLimit(ratelimit.ActionLogin), h.Login)          // 用户登录
		g.POST("/logout", mw.Auth(), h.Logout)                                  // 用户登出
		// g.GET("/me", mw.AuthMW(), h.Me)                   // 获取当前用户信息(该功能不适合放在 验证业务中)

		// 账号管理
		accountGroup := g.Group("/account", mw.Auth())
		{
			// DELETE /auth/account - 软删除账号
			accountGroup.DELETE("", h.DeleteAccount)
			// GET /auth/account/deletion - 获取删除状态
			accountGroup.GET("/deletion", h.DeletionStatus)
			// POST /auth/account/restore - 恢复账号
			accountGroup.POST("/restore", h.CancelDeletion)
			// DELETE /auth/account/permanent - 永久删除
			accountGroup.DELETE("/permanent", h.ConfirmDeletion)
			// PUT /auth/account/password - 修改密码
			accountGroup.PUT("/password", h.ChangePassword)

		}

		passwordGroup := g.Group("/password")
		{
			// POST /auth/password/forgot - 忘记密码
			passwordGroup.POST("/forgot", h.ForgotPassword) // 发送重置密码的邮件给用户
			// POST /auth/password/reset - 重置密码
			passwordGroup.PUT("/reset", h.ResetPassword) // 重置密码（存在安全问题，计划移除）
			// GET /auth/password/validate-token - 验证token
			passwordGroup.GET("/validate-token", h.ValidateResetToken)      // 用户点击邮件链接，验证 Token 是否有效
			passwordGroup.PUT("/reset-withtoken", h.ResetPasswordWithToken) // 通过 token 重置密码
			// GET /auth/password/reset
			// passwordGroup.GET("/reset", h.ShowResetPage)

		}
	}

}
