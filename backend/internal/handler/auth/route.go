package auth

// TODO: Refactory，不符合 REATful 规范
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
// base URL: /api/v1
// Group URL: /auth
func (h *AuthHandler) RegisterRoutes(auth *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	// 用户排行榜
	g := auth.Group("/auth")
	{
		g.POST("/register", h.Register)                               // 用户注册
		g.POST("/login", h.Login)                                     // 用户登录
		g.POST("/logout", h.Logout)                                   // 用户登出
// g.GET("/me", mw.AuthMW(), h.Me)                   // 获取当前用户信息(该功能不适合放在 验证业务中)

// 账号管理
accountGroup := g.Group("/account", mw.AuthMW())
		{
			// DELETE /auth/account - 软删除账号
			accountGroup.DELETE("", h.DeleteAccount)
			// GET /auth/account/deletion - 获取删除状态
			accountGroup.GET("/deletion", h.DeletionStatus)
			// POST /auth/account/restore - 恢复账号
			accountGroup.POST("/restore", h.CancelDeletion)
			// DELETE /auth/account/permanent - 永久删除
			accountGroup.DELETE("/permanent", h.ConfirmDeletion)
		}

		passwordGroup := g.Group("/password")
		{
			// POST /auth/password/forgot - 忘记密码
			passwordGroup.POST("/forgot", h.ForgotPassword)
			// POST /auth/password/reset - 重置密码
			passwordGroup.POST("/reset", h.ResetPassword)
			// GET /auth/password/validate-token - 验证token
			passwordGroup.GET("/validate-token", h.ValidateResetToken)
		}
		// === delete
		// g.DELETE("/account", mw.AuthMW(), h.DeleteAccount)     // 用户删除账号
		// g.GET("/deletion-status", mw.AuthMW(), h.DeletionStatus)      // 用户查询账号删除状态
		// g.POST("/cancel-deletion", mw.AuthMW(), h.CancelDeletion)     // 用户取消账号删除
		// g.DELETE("/confirm-deletion", mw.AuthMW(), h.ConfirmDeletion) // 用户确认账号删除
		
		// g.POST("/forgot-password", h.ForgotPassword)         // 用户忘记密码
		// g.POST("/reset-password", h.ResetPassword)           // 用户重置密码
		// g.GET("/validate-reset-token", h.ValidateResetToken) // 用户验证重置密码 token
	}

}
