package user

// TODO: Refactory
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes(user *gin.RouterGroup, mw *middleware.MiddlewareSet) {
	// 用户排行榜
	g := user.Group("/leaderboard")
	{
		g.GET("/simple", h.LeaderboardSimple) // 获取排行榜简要信息
		g.GET("/detail", h.LeaderboardDetail) // 获取排行榜详细信息

	}
	// 用户社交信息
	g = user.Group("/:id")
	{
		g.GET("", mw.OptionalAuthMW(), h.GetProfile)             // 获取用户信息（）
		g.POST("/follow", mw.AuthMW(), h.Follow)                 // 关注用户
		g.DELETE("/follow", mw.AuthMW(), h.Unfollow)             // 取消关注用户
		g.GET("/followers", mw.OptionalAuthMW(), h.GetFollowers) // 获取用户粉丝列表
		g.GET("/following", mw.OptionalAuthMW(), h.GetFollowing) // 获取用户关注列表
		g.GET("/Score", mw.OptionalAuthMW(), h.GetScore)         // 获取用户积分

	}
	// 用户个人信息
	g = user.Group("/me")
	{
		g.GET("/role", mw.OptionalAuthMW(), h.GetCurrentUserRole) // 获取当前用户角色
	}
	g = user.Group("/profile")
	{
		g.PUT("", mw.AuthMW(), h.UpdateProfile)

	}
	g = user.Group("/password")
	{
		g.PATCH("", mw.AuthMW(), h.ChangePassword)
	}
}
