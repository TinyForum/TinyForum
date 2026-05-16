package board

import (
	"strconv"

	"tiny-forum/internal/model/request"
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
// @Success 200 {object} common.BasicResponse  "添加版主成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误或板块ID无效"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（需要管理员或 manage_moderator 权限）"
// @Failure 404 {object} common.BasicResponse"板块不存在"
// @Router /boards/{id}/moderators [post]
// AddModerator 管理员直接添加版主
// POST /api/boards/:id/moderators
func (h *BoardHandler) AddModerator(c *gin.Context) {
	// 1. 解析并校验板块ID
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || boardID == 0 {
		response.HandleError(c, err)
		return
	}

	// 2. 绑定请求体（使用已定义的 DTO）
	var req request.AddModeratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	req.BoardID = uint(boardID)

	// 3. 请求参数校验（非空、权限合法性等）
	if req.UserID == 0 {
		response.HandleError(c, err)
		return
	}
	if err := req.Validate(); err != nil {
		response.HandleError(c, err)
		return
	}

	// 4. 获取操作人ID（从认证上下文）
	operatorID, err := getAuthenticatedUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	// 5. 构建 Service 层输入参数
	input := request.AddModeratorRequest{
		UserID:      req.UserID,
		BoardID:     req.BoardID,
		Permissions: req.Permissions,
	}

	// 6. 调用 Service
	if err := h.boardSvc.AddModerator(c.Request.Context(), input, operatorID); err != nil {
		response.HandleError(c, err)
		return
	}

	// 7. 成功响应
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
// @Success 200 {object} common.BasicResponse  "移除版主成功"
// @Failure 400 {object} common.BasicResponse"无效的板块ID或用户ID"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（需要管理员或 manage_moderator 权限）"
// @Failure 404 {object} common.BasicResponse"版主不存在"
// @Router /boards/{id}/moderators/{user_id} [delete]
func (h *BoardHandler) RemoveModerator(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	operatorID := c.GetUint("user_id")
	if err := h.boardSvc.RemoveModerator(c.Request.Context(), uint(userID), uint(boardID), operatorID); err != nil {
		response.HandleError(c, err)
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
// @Success 200 {object} common.BasicResponse  "获取成功"
// @Failure 400 {object} common.BasicResponse"无效的板块ID"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /boards/{id}/moderators [get]
func (h *BoardHandler) GetModerators(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	moderators, err := h.boardSvc.GetModerators(uint(boardID))
	if err != nil {
		response.HandleError(c, err)
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
// @Success 200 {object} common.BasicResponse  "权限更新成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误或板块ID/用户ID无效"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（需要管理员权限）"
// @Failure 404 {object} common.BasicResponse"版主不存在"
// @Router /boards/{id}/moderators/{user_id}/permissions [put]
func (h *BoardHandler) UpdateModeratorPermissions(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	var req request.UpdateModeratorPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	operatorID := c.GetUint("user_id")
	input := request.UpdateModeratorPermissionsRequest{
		UserID:      uint(userID),
		BoardID:     uint(boardID),
		Permissions: req.Permissions,
	}

	if err := h.boardSvc.UpdateModeratorPermissions(c.Request.Context(), input, operatorID); err != nil {
		response.HandleError(c, err)
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
// @Success 200 {object} common.BasicResponse "获取成功"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /boards/moderators/managed [get]
func (h *BoardHandler) GetUserModeratorBoards(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Unauthorized(c, "未登录")
		return
	}

	boards, err := h.boardSvc.GetModeratorBoardsWithPermissions(userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, boards)
}
