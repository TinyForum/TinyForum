package board

import (
	"tiny-forum/internal/middleware"
	"tiny-forum/internal/repository/board"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册版块相关路由
// api: 版块根路由
// - mw: 中间件集合
// - repo: 版块仓库
func (h *BoardHandler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet, repo board.BoardRepository) {
	// 版块根路由
	boards := api.Group("/boards")
	{
		// 公开路由
		boards.GET("", h.List)         // GET /api/v1/boards  获取版块列表
		boards.GET("/tree", h.GetTree) // GET /api/v1/boards/tree 获取版块树

		// 通过 slug 访问（公开）
		boards.GET("/slug/:slug", h.GetBoardBySlug)                          // GET /api/v1/boards/slug/:slug  获取版块信息
		boards.GET("/slug/:slug/posts", mw.OptionalAuth(), h.GetPostsBySlug) // GET /api/v1/boards/slug/:slug/posts  获取版块帖子列表

		// 通过 ID 访问（公开）
		boards.GET("/:id", mw.OptionalAuth(), h.GetByID) // GET /api/v1/boards/:id 获取版块信息

		// 版主申请（需要认证）
		boards.POST("/:id/moderators/apply", mw.Auth(), h.ApplyModerator)              // POST /api/v1/boards/:id/moderators/apply 申请成为版主
		boards.GET("/moderators/apply-status", mw.Auth(), h.GetUserApplications)       // GET /api/v1/boards/moderators/apply-status 获取用户的版主申请状态
		boards.GET("/moderators/managed", mw.Auth(), h.GetUserModeratorBoards)         // GET /api/v1/boards/moderators/managed 获取用户管理的版块列表
		boards.DELETE("/applications/:application_id", mw.Auth(), h.CancelApplication) // DELETE /api/v1/boards/applications/:application_id 取消版主申请

		// 管理员操作
		adminGroup := boards.Group("")
		adminGroup.Use(mw.Auth(), mw.AdminRequired())
		{
			adminGroup.POST("", h.Create)       // POST /api/v1/boards 创建版块
			adminGroup.PUT("/:id", h.Update)    // PUT /api/v1/boards/:id 更新版块
			adminGroup.DELETE("/:id", h.Delete) // DELETE /api/v1/boards/:id 删除版块
		}

		// 版主管理（需要版主权限）
		moderatorGroup := boards.Group("/:id/moderators")
		moderatorGroup.Use(mw.Auth())
		{
			moderatorGroup.GET("", mw.ModeratorRequired(repo), h.GetModerators)                           // GET /api/v1/boards/:id/moderators 获取版主列表
			moderatorGroup.POST("", mw.CanManageModerator(repo), h.AddModerator)                          // POST /api/v1/boards/:id/moderators 添加版主
			moderatorGroup.DELETE("/:user_id", mw.CanManageModerator(repo), h.RemoveModerator)            // DELETE /api/v1/boards/:id/moderators/:user_id 移除版主
			moderatorGroup.PUT("/:user_id/permissions", mw.AdminRequired(), h.UpdateModeratorPermissions) // PUT /api/v1/boards/:id/moderators/:user_id/permissions 更新版主权限
		}

		// 封禁管理（需要封禁权限）
		banGroup := boards.Group("/:id/bans")
		banGroup.Use(mw.Auth())
		{
			banGroup.POST("", mw.CanBanUser(repo), h.BanUser)              // POST /api/v1/boards/:id/bans 封禁用户
			banGroup.DELETE("/:user_id", mw.CanBanUser(repo), h.UnbanUser) // DELETE /api/v1/boards/:id/bans/:user_id 解封用户
		}

		// 帖子管理（需要管理权限）
		postManageGroup := boards.Group("/:id/posts")
		postManageGroup.Use(mw.Auth())
		{
			postManageGroup.DELETE("/:post_id", mw.CanDeletePost(repo), h.DeletePost) // DELETE /api/v1/boards/:id/posts/:post_id 删除帖子
			postManageGroup.PUT("/:post_id/pin", mw.CanPinPost(repo), h.PinPost)      // PUT /api/v1/boards/:id/posts/:post_id/pin 置顶帖子
		}
	}
}
