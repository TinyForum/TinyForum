package middleware

import (
	"strconv"

	"tiny-forum/internal/repository"
	"tiny-forum/pkg/jwt"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ModeratorRequired 检查用户是否为版主（任意权限即可）
func ModeratorRequired(jwtMgr *jwt.Manager, boardRepo *repository.BoardRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 检查是否为管理员（管理员拥有所有版主权限）
		userRole, exists := c.Get("user_role")
		if exists && userRole == "admin" || exists && userRole == "super_admin" {
			c.Next()
			return
		}

		// 获取板块ID
		boardIDStr := c.Param("id")
		if boardIDStr == "" {
			boardIDStr = c.Param("board_id")
		}
		if boardIDStr == "" {
			response.BadRequest(c, "缺少板块ID")
			c.Abort()
			return
		}

		boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的板块ID")
			c.Abort()
			return
		}

		// 检查是否为版主
		isMod, err := boardRepo.IsModerator(userID.(uint), uint(boardID))
		if err != nil {
			response.InternalError(c, "检查权限失败")
			c.Abort()
			return
		}

		if !isMod {
			response.Forbidden(c, "需要版主权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

// SpecificModeratorPermission 检查版主是否有特定权限
func SpecificModeratorPermission(jwtMgr *jwt.Manager, boardRepo *repository.BoardRepository, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		// 管理员拥有所有权限
		userRole, exists := c.Get("user_role")
		if exists && userRole == "admin" {
			c.Next()
			return
		}

		// 获取板块ID
		boardIDStr := c.Param("id")
		if boardIDStr == "" {
			boardIDStr = c.Param("board_id")
		}
		if boardIDStr == "" {
			response.BadRequest(c, "缺少板块ID")
			c.Abort()
			return
		}

		boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的板块ID")
			c.Abort()
			return
		}

		// 获取版主信息
		moderator, err := boardRepo.FindModeratorByUserAndBoard(userID.(uint), uint(boardID))
		if err != nil {
			response.Forbidden(c, "不是版主")
			c.Abort()
			return
		}

		// 检查特定权限
		// hasPermission := false
		hasPermission := moderator.HasPermission(permission)

		if !hasPermission {
			response.Forbidden(c, "没有权限执行此操作")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CanDeletePost 检查是否有删除帖子的权限
func CanDeletePost(jwtMgr *jwt.Manager, boardRepo *repository.BoardRepository) gin.HandlerFunc {
	return SpecificModeratorPermission(jwtMgr, boardRepo, "delete_post")
}

// CanPinPost 检查是否有置顶帖子的权限
func CanPinPost(jwtMgr *jwt.Manager, boardRepo *repository.BoardRepository) gin.HandlerFunc {
	return SpecificModeratorPermission(jwtMgr, boardRepo, "pin_post")
}

// CanEditAnyPost 检查是否有编辑任何帖子的权限
func CanEditAnyPost(jwtMgr *jwt.Manager, boardRepo *repository.BoardRepository) gin.HandlerFunc {
	return SpecificModeratorPermission(jwtMgr, boardRepo, "edit_any_post")
}

// CanManageModerator 检查是否有管理版主的权限
func CanManageModerator(jwtMgr *jwt.Manager, boardRepo *repository.BoardRepository) gin.HandlerFunc {
	return SpecificModeratorPermission(jwtMgr, boardRepo, "manage_moderator")
}

// CanBanUser 检查是否有禁言用户的权限
func CanBanUser(jwtMgr *jwt.Manager, boardRepo *repository.BoardRepository) gin.HandlerFunc {
	return SpecificModeratorPermission(jwtMgr, boardRepo, "ban_user")
}
