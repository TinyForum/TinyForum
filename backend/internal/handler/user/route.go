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
		leaderboard.GET("/simple", h.LeaderboardSimple) // GET /api/v1/users/leaderboard/simple 获取排行榜简单信息
		leaderboard.GET("/detail", h.LeaderboardDetail) // GET /api/v1/users/leaderboard/detail 获取排行榜详细信息
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
		stats := me.Group("/stats")
		{

			stats.GET("", h.GetStatisticsCount) // GET /api/v1/users/me/stats 用户统计信息

		}

		posts := me.Group("/posts")
		{
			posts.GET("", h.GetUserPosts)              // GET /api/v1/users/me/posts 用户帖子列表
			posts.GET("/:post_id", h.GetUserPostById)  // GET /api/v1/users/me/posts/:post_id 我的单个帖子
		}
		me.GET("/comments", h.GetUserComments) // GET /api/v1/users/me/comments 我的评论列表
		violation := me.Group("/violations")
		{
			violation.GET("", h.ListUserViolations)              // GET /api/v1/users/me/violations 用户违规记录
			violation.GET("/:id", h.GetUserViolation)            // GET /api/v1/users/me/violations/:id 违规详情
			violation.POST("/:id/appeal", h.AppealUserViolation) // POST /api/v1/users/me/violations/:id/appeal 违规申诉
		}
		// me.PATCH("/password", h.ChangePassword) // PATCH /api/v1/users/me/password 修改密码
	}

	// 4. 管理员用户管理（密码重置、积分设置）
	adminUsers := api.Group("/admin/users", mw.Auth(), mw.AdminRequired())
	{
		adminUsers.POST("/:id/reset-password", h.AdminResetUserPassword) // POST /api/v1/admin/users/:id/reset-password 重置用户密码
		adminUsers.PUT("/:id/score", h.AdminSetScore)                    // PUT /api/v1/admin/users/:id/score 设置用户积分
	}

	// 5. 当前登录用户信息（前端登录态刷新使用）
	authGroup := api.Group("/auth")
	{
		authGroup.GET("/me", mw.Auth(), h.Me) // GET /api/v1/auth/me 获取当前用户信息
	}
}
