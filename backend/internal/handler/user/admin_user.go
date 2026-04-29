package user

import (
	"strconv"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminList 管理员获取用户列表
// @Summary 管理员获取用户列表
// @Tags 管理接口
// @Produce json
// @Security ApiKeyAuth
// @Router /admin/users [get]
func (h *UserHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	users, total, err := h.userSvc.List(page, pageSize, keyword)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, users, total, page, pageSize)
}

// @Summary 管理员设置用户状态
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param body body object{is_active=bool} true "状态"
// @Router /admin/users/{id}/active [put]
func (h *UserHandler) AdminSetActive(c *gin.Context) {
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
	if err := h.userSvc.SetActive(uint(targetID), currentID, body.IsActive); err != nil {
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
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param body body object{is_blocked=bool} true "封禁状态"
// @Router /admin/users/{id}/blocked [put]
func (h *UserHandler) AdminSetBlocked(c *gin.Context) {
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
	if err := h.userSvc.SetBlocked(uint(targetID), currentID, body.IsBlocked); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{
		"user_id":    targetID,
		"is_blocked": body.IsBlocked,
	})
}

// AdminSetRole 设置用户角色
// @Summary 管理员设置用户角色
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param body body SetUserRoleRequest true "角色信息"
// @Success 200 {object} response.Response
// @Router /admin/users/{id}/role [put]
func (h *UserHandler) AdminSetRole(c *gin.Context) {
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
	var body SetUserRoleRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}
	if err := h.userSvc.SetRole(operatorID.(uint), uint(targetID), body.Role); err != nil {
		handleRoleError(c, err)
		return
	}
	response.Success(c, gin.H{
		"message":     "设置角色成功",
		"user_id":     targetID,
		"new_role":    body.Role,
		"operator_id": operatorID,
	})
}

// AdminDeleteUser 管理员删除用户
// @Summary 管理员删除用户
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "目标用户ID"
// @Router /admin/users/{id} [delete]
func (h *UserHandler) AdminDeleteUser(c *gin.Context) {
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
	if err := h.userSvc.DeleteUser(operatorUint, uint(targetID)); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{
		"message":     "删除用户成功",
		"user_id":     targetID,
		"operator_id": operatorUint,
	})
}

// AdminResetUserPassword 管理员重置用户密码
// @Summary 管理员重置用户密码（生成临时密码并通知用户）
// @Description 管理员为指定用户生成随机临时密码，并通过站内通知发送给用户。
// @Description 临时密码有效期为 30 分钟，用户登录后需尽快修改密码。
// @Description
// @Description **权限要求**：
// @Description - 超级管理员：可重置任何用户
// @Description - 普通管理员：只能重置普通用户，不能重置管理员和超级管理员
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "目标用户ID" example:"123"
// @Success 200 {object} response.Response{data=AdminResetUserPasswordResponse} "操作成功"
// @Failure 400 {object} response.Response{data=[]response.ValidationError} "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "权限不足"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/users/{id}/reset-password [post]
func (h *UserHandler) AdminResetUserPassword(c *gin.Context) {
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
	tempPassword, err := h.userSvc.ResetUserPasswordWithTemp(operatorUint, uint(targetID))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	h.sendTempPasswordNotification(uint(targetID), operatorUint, tempPassword)
	response.Success(c, AdminResetUserPasswordResponse{
		Message:    "临时密码已生成并发送给用户",
		UserID:     uint(targetID),
		OperatorID: operatorUint,
	})
}
