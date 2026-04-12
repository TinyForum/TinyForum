package handler

import (
	"errors"
	"strconv"

	apperrors "tiny-forum/internal/errors"
	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// ── 公开接口 ─────────────────────────────────────────────────────────────────

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=service.UserProfileResponse}
// @Router /users/{id} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	viewerUint, _ := c.Get("user_id")
	var viewerID uint
	if v, ok := viewerUint.(uint); ok {
		viewerID = v
	}
	profile, err := h.userSvc.GetUserProfile(uint(targetID), viewerID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, profile)
}

// UpdateProfile 更新个人资料
// @Summary 更新用户资料
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.UpdateProfileInput true "资料"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input service.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.userSvc.UpdateProfile(userID, input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	user, _ := h.userSvc.GetProfile(userID)
	response.Success(c, user)
}

// Follow 关注用户
// @Summary 关注用户
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Router /users/{id}/follow [post]
func (h *UserHandler) Follow(c *gin.Context) {
	followerID := c.GetUint("user_id")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	if err := h.userSvc.Follow(followerID, uint(targetID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "关注成功"})
}

// Unfollow 取消关注
// @Summary 取消关注用户
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Router /users/{id}/follow [delete]
func (h *UserHandler) Unfollow(c *gin.Context) {
	followerID := c.GetUint("user_id")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	if err := h.userSvc.Unfollow(followerID, uint(targetID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "已取消关注"})
}

// Leaderboard 用户排行榜
// @Summary 获取用户排行榜
// @Tags 用户管理
// @Produce json
// @Param limit query int false "数量" default(20)
// @Router /users/leaderboard [get]
func (h *UserHandler) Leaderboard(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	users, err := h.userSvc.GetLeaderboard(limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, users)
}

// GetFollowers 获取粉丝列表
// @Summary 获取粉丝列表
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Router /users/{id}/followers [get]
func (h *UserHandler) GetFollowers(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	followers, total, err := h.userSvc.GetFollowers(uint(userID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, followers, total, page, pageSize)
}

// GetFollowing 获取关注列表
// @Summary 获取关注列表
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Router /users/{id}/following [get]
func (h *UserHandler) GetFollowing(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	following, total, err := h.userSvc.GetFollowing(uint(userID), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, following, total, page, pageSize)
}

// ── 管理接口 ─────────────────────────────────────────────────────────────────

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

// AdminSetActive 设置用户启用状态
// @Summary 管理员设置用户状态
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Router /admin/users/{id}/active [put]
func (h *UserHandler) AdminSetActive(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	var body struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.userSvc.SetActive(uint(targetID), body.IsActive); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "操作成功"})
}

// AdminSetBlocked 设置用户封禁状态
// @Summary 管理员封禁/解封用户
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Router /admin/users/{id}/blocked [put]
func (h *UserHandler) AdminSetBlocked(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	var body struct {
		IsBlocked bool `json:"is_blocked"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.userSvc.SetBlocked(uint(targetID), body.IsBlocked); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "操作成功"})
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
		h.handleRoleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"message":     "设置角色成功",
		"user_id":     targetID,
		"new_role":    body.Role,
		"operator_id": operatorID,
	})
}

// handleRoleError 统一处理角色变更错误，集中映射到 HTTP 响应
func (h *UserHandler) handleRoleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		response.NotFound(c, "用户不存在")
	case errors.Is(err, apperrors.ErrInvalidRole):
		response.BadRequest(c, "无效的角色类型")
	case errors.Is(err, apperrors.ErrCannotModifySelf):
		response.Forbidden(c, "不能修改自己的角色")
	case errors.Is(err, apperrors.ErrCannotChangeOwnerRole):
		response.Forbidden(c, "不能修改超级管理员的角色")
	case errors.Is(err, apperrors.ErrInsufficientPermission):
		response.Forbidden(c, "权限不足："+err.Error())
	default:
		response.InternalError(c, "设置角色失败: "+err.Error())
	}
}

// ── 请求体结构 ────────────────────────────────────────────────────────────────

// SetUserRoleRequest 设置用户角色请求
type SetUserRoleRequest struct {
	// 角色：user / member / moderator / reviewer / bot / admin / super_admin
	Role string `json:"role" binding:"required,oneof=user member moderator reviewer bot admin super_admin"`
}

// SetUserActiveRequest 设置用户激活状态请求
type SetUserActiveRequest struct {
	IsActive bool `json:"is_active" example:"true"`
}

// SetUserBlockedRequest 设置用户封禁状态请求
type SetUserBlockedRequest struct {
	IsBlocked bool `json:"is_blocked" example:"true"`
}
