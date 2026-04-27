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
	g := auth.Group("")
	{
		g.POST("/register", h.Register)                               // 用户注册
		g.POST("/login", h.Login)                                     // 用户登录
		g.POST("/logout", h.Logout)                                   // 用户登出
		g.DELETE("/delete-account", mw.AuthMW(), h.DeleteAccount)     // 用户删除账号
		g.GET("/deletion-status", mw.AuthMW(), h.DeletionStatus)      // 用户查询账号删除状态
		g.POST("/cancel-deletion", mw.AuthMW(), h.CancelDeletion)     // 用户取消账号删除
		g.DELETE("/confirm-deletion", mw.AuthMW(), h.ConfirmDeletion) // 用户确认账号删除
		// g.GET("/me", mw.AuthMW(), handlers.User.Me)                   // 获取当前用户信息
		g.POST("/forgot-password", h.ForgotPassword)         // 用户忘记密码
		g.POST("/reset-password", h.ResetPassword)           // 用户重置密码
		g.GET("/validate-reset-token", h.ValidateResetToken) // 用户验证重置密码 token
	}

}
