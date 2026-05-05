package admin

import (
	"strconv"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminList 管理员获取用户列表
// @Summary 管理员获取用户列表
// @Tags 管理员后台
// @Produce json
// @Security ApiKeyAuth
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	users, total, err := h.service.ListUsers(page, pageSize, keyword)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, users, total, page, pageSize)
}

// @Summary 管理员设置用户状态
// @Tags 管理员后台
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param body body object{is_active=bool} true "状态"
// @Router /admin/users/{id}/active [put]
func (h *AdminHandler) SetActiveUser(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFailed(c, []response.ValidationError{
			{Field: "id", Message: "无效的用户ID格式"},
		})
		return
	}
	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	currentID, ok := currentUserID.(uint)
	if !ok {
		response.InternalError(c, "用户身份解析失败")
		return
	}
	if currentID == uint(targetID) {
		response.Forbidden(c, "不能修改自己的激活状态")
		return
	}
	var body struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ValidationFailed(c, []response.ValidationError{
			{Field: "is_active", Message: "is_active 字段必须为布尔值"},
		})
		return
	}
	if err := h.service.SetActiveUser(uint(targetID), currentID, body.IsActive); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{
		"user_id":   targetID,
		"is_active": body.IsActive,
	})
}

// AdminSetBlocked 设置用户封禁状态
// @Summary 管理员封禁/解封用户
// @Tags 管理员后台
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param body body object{is_blocked=bool} true "封禁状态"
// @Router /admin/users/{id}/blocked [put]
func (h *AdminHandler) SetBlockedUser(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFailed(c, response.SimpleValidationError("id", "无效的用户ID格式"))
		return
	}
	currentUserID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	currentID, ok := currentUserID.(uint)
	if !ok {
		response.InternalError(c, "用户身份解析失败")
		return
	}
	var body struct {
		IsBlocked bool `json:"is_blocked"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.ValidationFailed(c, response.SimpleValidationError("is_blocked", "is_blocked 字段必须为布尔值"))
		return
	}
	if err := h.service.SetBlockedUser(uint(targetID), currentID, body.IsBlocked); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{
		"user_id":    targetID,
		"is_blocked": body.IsBlocked,
	})
}

// AdminDeleteUser 管理员删除用户
// @Summary 管理员删除用户
// @Tags 管理员后台
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "目标用户ID"
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权操作")
		return
	}
	operatorUint, ok := operatorID.(uint)
	if !ok {
		response.InternalError(c, "用户身份解析失败")
		return
	}
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationFailed(c, response.SimpleValidationError("id", "无效的用户ID格式"))
		return
	}
	if err := h.service.DeleteUser(operatorUint, uint(targetID)); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{
		"message":     "删除用户成功",
		"user_id":     targetID,
		"operator_id": operatorUint,
	})
}

// AdminSetRole 设置用户角色
// @Summary 管理员设置用户角色
// @Tags 管理员后台
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param body body request.SetUserRoleRequest true "角色信息"
// @Success 200 {object} vo.BasicResponse
// @Router /admin/users/{id}/role [put]
func (h *AdminHandler) SetRoleUser(c *gin.Context) {
	operatorID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权操作")
		return
	}
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	var body request.SetUserRoleRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.Errorf("修改用户角色失败: ",err.Error())
		response.BadRequest(c, "请求参数错误")
		return
	}
	if err := h.service.SetRoleUser(operatorID.(uint), uint(targetID), body.Role); err != nil {
		logger.Errorf("修改用户角色失败: ",err.Error())
		response.HandleError(c, err)
		return
	}
	resp := vo.AdminSetUserRole{
	Message:    "设置角色成功",
	UserID:     targetID,
	NewRole:    body.Role,
	OperatorID: operatorID.(uint),
}
response.Success(c, resp)
}
