package user

// TODO: Refactory
import (
	"tiny-forum/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func (h *UserHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	// 用户资源根路径
	users := api.Group("/users")

	// 1. 排行榜（独立资源）
	leaderboard := users.Group("/leaderboard")
	{
		leaderboard.GET("/simple", h.LeaderboardSimple) // GET /api/v1/users/leaderboard/simple
		leaderboard.GET("/detail", h.LeaderboardDetail) // GET /api/v1/users/leaderboard/detail
	}

	// 2. 特定用户的操作
	user := users.Group("/:id")
	{
		// 公开/可选认证
		user.GET("", mw.OptionalAuth(), h.GetProfile)             // GET /api/v1/users/:id  获取用户信息
		user.GET("/followers", mw.OptionalAuth(), h.GetFollowers) // GET /api/v1/users/:id/followers 获取关注者列表
		user.GET("/following", mw.OptionalAuth(), h.GetFollowing) // GET /api/v1/users/:id/following 获取关注列表
		user.GET("/score", mw.OptionalAuth(), h.GetScore)         // GET /api/v1/users/:id/score 获取用户积分

		// 需要认证
		auth := user.Group("")
		auth.Use(mw.Auth())
		{
			auth.POST("/follow", h.Follow)     // POST /api/v1/users/:id/follow 关注用户
			auth.DELETE("/follow", h.Unfollow) // DELETE /api/v1/users/:id/follow 取消关注
		}
	}

	// 3. 当前用户自己的信息
	me := users.Group("/me")
	me.Use(mw.Auth()) // 所有 /me 操作都需要认证
	{
		me.GET("/role", h.GetCurrentUserRole) // GET /api/v1/users/me/role 获取当前用户角色
		me.PUT("/profile", h.UpdateProfile)   // PUT /api/v1/users/me/profile 更新用户信息
		// me.PATCH("/password", h.ChangePassword) // PATCH /api/v1/users/me/password 修改密码
	}
}
