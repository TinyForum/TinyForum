package board

import (
	"strconv"

	boardService "tiny-forum/internal/service/board"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// AddModerator 直接任命版主（管理员 / 有 manage_moderator 权限的版主）
// @Summary 任命版主
// @Description 管理员或有 manage_moderator 权限的版主可直接任命版主并设置权限
// @Tags 版主管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param body body object true "版主信息"
// @Param body.user_id body int true "用户ID" example(10086)
// @Param body.can_delete_post body bool false "删除帖子权限" example(true)
// @Param body.can_pin_post body bool false "置顶帖子权限" example(true)
// @Param body.can_edit_any_post body bool false "编辑任意帖子权限" example(false)
// @Param body.can_manage_moderator body bool false "管理版主权限" example(false)
// @Param body.can_ban_user body bool false "禁言用户权限" example(true)
// @Success 200 {object} response.Response{data=object} "添加版主成功"
// @Failure 400 {object} response.Response "请求参数错误或板块ID无效"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（需要管理员或 manage_moderator 权限）"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/{id}/moderators [post]
func (h *BoardHandler) AddModerator(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}

	var body struct {
		UserID             uint `json:"user_id"              binding:"required"`
		CanDeletePost      bool `json:"can_delete_post"`
		CanPinPost         bool `json:"can_pin_post"`
		CanEditAnyPost     bool `json:"can_edit_any_post"`
		CanManageModerator bool `json:"can_manage_moderator"`
		CanBanUser         bool `json:"can_ban_user"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	input := boardService.AddModeratorInput{
		UserID:             body.UserID,
		BoardID:            uint(boardID),
		CanDeletePost:      body.CanDeletePost,
		CanPinPost:         body.CanPinPost,
		CanEditAnyPost:     body.CanEditAnyPost,
		CanManageModerator: body.CanManageModerator,
		CanBanUser:         body.CanBanUser,
	}

	if err := h.boardSvc.AddModerator(c.Request.Context(), input, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "添加版主成功"})
}

// RemoveModerator 移除版主
// @Summary 移除版主
// @Description 管理员或有 manage_moderator 权限的版主可移除指定板块的版主
// @Tags 版主管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param user_id path int true "用户ID"
// @Success 200 {object} response.Response{data=object} "移除版主成功"
// @Failure 400 {object} response.Response "无效的板块ID或用户ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（需要管理员或 manage_moderator 权限）"
// @Failure 404 {object} response.Response "版主不存在"
// @Router /boards/{id}/moderators/{user_id} [delete]
func (h *BoardHandler) RemoveModerator(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.boardSvc.RemoveModerator(c.Request.Context(), uint(userID), uint(boardID), operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "移除版主成功"})
}

// GetModerators 获取板块版主列表
// @Summary 获取板块版主列表
// @Description 获取指定板块的所有版主信息
// @Tags 版主管理
// @Produce json
// @Param id path int true "板块ID"
// @Success 200 {object} response.Response{data=[]model.Moderator} "获取成功"
// @Failure 400 {object} response.Response "无效的板块ID"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /boards/{id}/moderators [get]
func (h *BoardHandler) GetModerators(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	moderators, err := h.boardSvc.GetModerators(uint(boardID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, moderators)
}

// UpdateModeratorPermissions 升级/降级版主权限（管理员）
// @Summary 更新版主权限
// @Description 管理员更新指定版主的权限配置
// @Tags 版主管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param user_id path int true "用户ID"
// @Param body body object true "权限配置"
// @Param body.can_delete_post body bool false "删除帖子权限" example(true)
// @Param body.can_pin_post body bool false "置顶帖子权限" example(true)
// @Param body.can_edit_any_post body bool false "编辑任意帖子权限" example(false)
// @Param body.can_manage_moderator body bool false "管理版主权限" example(false)
// @Param body.can_ban_user body bool false "禁言用户权限" example(true)
// @Success 200 {object} response.Response{data=object} "权限更新成功"
// @Failure 400 {object} response.Response "请求参数错误或板块ID/用户ID无效"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（需要管理员权限）"
// @Failure 404 {object} response.Response "版主不存在"
// @Router /boards/{id}/moderators/{user_id}/permissions [put]
func (h *BoardHandler) UpdateModeratorPermissions(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	var body struct {
		CanDeletePost      bool `json:"can_delete_post"`
		CanPinPost         bool `json:"can_pin_post"`
		CanEditAnyPost     bool `json:"can_edit_any_post"`
		CanManageModerator bool `json:"can_manage_moderator"`
		CanBanUser         bool `json:"can_ban_user"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	input := boardService.UpdateModeratorPermissionsInput{
		UserID:             uint(userID),
		BoardID:            uint(boardID),
		CanDeletePost:      body.CanDeletePost,
		CanPinPost:         body.CanPinPost,
		CanEditAnyPost:     body.CanEditAnyPost,
		CanManageModerator: body.CanManageModerator,
		CanBanUser:         body.CanBanUser,
	}

	if err := h.boardSvc.UpdateModeratorPermissions(c.Request.Context(), input, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "权限更新成功"})
}

// GetUserModeratorBoards 获取当前用户管理的板块列表（含权限）
// @Summary 获取我管理的板块
// @Description 获取当前登录用户作为版主管理的所有板块，包含权限信息
// @Tags 版主管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]board.ModeratorBoardWithPerms} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /boards/moderators/managed [get]
func (h *BoardHandler) GetUserModeratorBoards(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	boards, err := h.boardSvc.GetModeratorBoardsWithPermissions(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, boards)
}
