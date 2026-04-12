package handler

import (
	"strconv"

	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 根据用户ID获取用户详细资料
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=service.UserProfileResponse} "获取成功"
// @Failure 400 {object} response.Response "无效的用户ID"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /users/{id} [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	viewerID, _ := c.Get("user_id")
	var viewerUint uint
	if v, ok := viewerID.(uint); ok {
		viewerUint = v
	}

	profile, err := h.userSvc.GetUserProfile(uint(targetID), viewerUint)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, profile)
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前登录用户的个人资料
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.UpdateProfileInput true "用户资料信息"
// @Success 200 {object} response.Response{data=service.UserProfileResponse} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
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
// @Description 关注指定用户
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "要关注的用户ID"
// @Success 200 {object} response.Response{data=object} "关注成功"
// @Failure 400 {object} response.Response "无效的用户ID或不能关注自己"
// @Failure 401 {object} response.Response "未授权"
// @Failure 409 {object} response.Response "已关注该用户"
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

// Unfollow 取消关注用户
// @Summary 取消关注用户
// @Description 取消关注指定用户
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "要取消关注的用户ID"
// @Success 200 {object} response.Response{data=object} "取消关注成功"
// @Failure 400 {object} response.Response "无效的用户ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "未关注该用户"
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

// Leaderboard 获取用户排行榜
// @Summary 获取用户排行榜
// @Description 根据用户积分或活跃度获取排行榜
// @Tags 用户管理
// @Produce json
// @Param limit query int false "返回数量" default(20) maximum(100)
// @Success 200 {object} response.Response{data=[]model.User} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /users/leaderboard [get]
func (h *UserHandler) Leaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
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

// AdminList 管理员获取用户列表
// @Summary 管理员获取用户列表
// @Description 管理员分页获取所有用户列表，支持关键词搜索
// @Tags 管理接口
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词（用户名、邮箱）"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.User}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 500 {object} response.Response "服务器内部错误"
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

// AdminSetActive 管理员设置用户状态
// @Summary 管理员设置用户状态
// @Description 管理员启用或禁用指定用户账号
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param body body SetUserActiveRequest true "状态信息"
// @Success 200 {object} response.Response{data=object} "操作成功"
// @Failure 400 {object} response.Response "无效的用户ID或请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "用户不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/users/{id}/active [put]
func (h *UserHandler) AdminSetActive(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	var body struct {
		Active bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.userSvc.SetActive(uint(targetID), body.Active); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "操作成功"})
}

// GetFollowers 获取粉丝列表
// @Summary 获取粉丝列表
// @Description 获取指定用户的粉丝列表
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.User}} "获取成功"
// @Failure 400 {object} response.Response "无效的用户ID"
// @Failure 500 {object} response.Response "服务器内部错误"
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
// @Description 获取指定用户关注的用户列表
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.User}} "获取成功"
// @Failure 400 {object} response.Response "无效的用户ID"
// @Failure 500 {object} response.Response "服务器内部错误"
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

// SetUserRoleRequest 设置用户角色请求参数
type SetUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user moderator admin" example:"moderator"` // 用户角色
}

// AdminSetRole 管理员设置用户角色
// @Summary 管理员设置用户角色
// @Description 管理员设置指定用户的角色（user/moderator/admin）
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param body body SetUserRoleRequest true "角色信息"
// @Success 200 {object} response.Response{data=object} "操作成功"
// @Failure 400 {object} response.Response "无效的用户ID或角色"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /admin/users/{id}/role [put]
func (h *UserHandler) AdminSetRole(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	var body SetUserRoleRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.userSvc.SetRole(uint(targetID), body.Role); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "设置角色成功"})
}

// SetUserActiveRequest 设置用户状态请求参数
type SetUserActiveRequest struct {
	Active bool `json:"active" example:"true"` // 用户状态：true-启用，false-禁用
}
