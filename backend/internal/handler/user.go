package handler

import (
	"errors"
	"strconv"
	"time"

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

// GetScore 查询用户积分（用户查自己，管理员可查任意用户）
// @Summary 查询积分
// @Tags 用户管理
// @Produce json
// @Security ApiKeyAuth
// @Router /users/score [get]
func (h *UserHandler) GetScore(c *gin.Context) {
	// 获取当前登录用户信息
	viewerID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权访问")
		return
	}

	viewerUint, ok := viewerID.(uint)
	if !ok {
		response.BadRequest(c, "无效的用户身份信息")
		return
	}

	// 获取用户角色
	viewerRole, _ := c.Get("role")
	viewerRoleStr, _ := viewerRole.(string)

	// 获取要查询的用户ID
	var targetID uint
	idParam := c.Param("id")

	if idParam != "" {
		// 有指定用户ID
		userID, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的用户ID")
			return
		}
		targetID = uint(userID)

		// 权限检查：普通用户只能查自己，管理员可以查所有
		if viewerUint != targetID && viewerRoleStr != "admin" && viewerRoleStr != "super_admin" {
			response.Forbidden(c, "权限不足，只能查询自己的积分")
			return
		}
	} else {
		// 没有指定用户ID，默认查询自己的
		targetID = viewerUint
	}

	// 查询积分
	score, err := h.userSvc.GetScoreById(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "用户不存在")
			return
		}
		response.InternalError(c, "查询积分失败: "+err.Error())
		return
	}

	// 返回成功响应
	response.Success(c, gin.H{
		"score":   score,
		"user_id": targetID,
	})
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

// AdminSetScore 管理员设置用户积分
// @Summary 管理员设置用户积分
// @Description 管理员可以通过此接口对指定用户的积分进行设置、增加或扣除操作。支持三种操作模式：<br>
// @Description - **set**：将用户积分设置为指定值<br>
// @Description - **add**：在现有积分基础上增加指定分数<br>
// @Description - **subtract**：从现有积分中扣除指定分数<br>
// @Description 积分范围限制为 0 ~ 999999，且操作后积分不能为负数或超出上限。
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID" example(10086)
// @Param body body AdminSetScoreRequest true "积分操作信息"
// @Success 200 {object} response.Response{data=AdminSetScoreResponse} "操作成功"
// @Failure 400 {object} response.Response "请求参数错误（如积分范围非法、操作类型错误等）"
// @Failure 401 {object} response.Response "未授权（缺少或无效的认证令牌）"
// @Failure 403 {object} response.Response "禁止访问（当前管理员无权限操作该用户）"
// @Failure 500 {object} response.Response "服务器内部错误（如数据库操作失败）"
// @Router /admin/users/{id}/score [put]
func (h *UserHandler) AdminSetScore(c *gin.Context) {
	// 1. 解析目标用户ID
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidUserID.Error())
		return
	}

	// 2. 解析请求体（使用已定义的结构体）
	var req AdminSetScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}

	// 3. 获取当前积分
	currentScore, err := h.userSvc.GetScoreById(uint(targetID))
	if err != nil {
		response.InternalError(c, "查询当前积分失败")
		return
	}

	// 4. 计算新积分
	var newScore int
	switch req.Operation {
	case "set":
		newScore = req.Score
	case "add":
		newScore = currentScore + req.Score
	case "subtract":
		newScore = currentScore - req.Score
	}

	// 5. 验证积分范围（防止负数或过大）
	if newScore < 0 {
		response.BadRequest(c, "积分不能为负数")
		return
	}
	if newScore > 999999 {
		response.BadRequest(c, "积分超出最大限制")
		return
	}

	// 6. 获取操作人信息
	viewerID, _ := c.Get("user_id")
	viewerUint, _ := viewerID.(uint)

	// 7. 设置新积分
	err = h.userSvc.SetScoreById(uint(targetID), newScore)
	if err != nil {
		response.InternalError(c, "设置积分失败: "+err.Error())
		return
	}

	// 8. 返回成功响应
	response.Success(c, AdminSetScoreResponse{
		UserID:     targetID,
		OldScore:   currentScore,
		NewScore:   newScore,
		Change:     newScore - currentScore,
		Operation:  req.Operation,
		OperatorID: viewerUint,
		Reason:     req.Reason,
		Timestamp:  time.Now().Unix(),
	})
}

// AdminSetScoreRequest 管理员设置积分请求参数
// swagger:model
type AdminSetScoreRequest struct {
	// 操作类型：set（设置）、add（增加）、subtract（扣除）
	// required: true
	// enum: set,add,subtract
	Operation string `json:"operation" binding:"required,oneof=set add subtract" example:"add"`

	// 积分数量（set 时为目标积分，add/subtract 时为变化量）
	// required: true
	// minimum: 0
	// maximum: 999999
	Score int `json:"score" binding:"required,gte=0,lte=999999" example:"50"`

	// 操作原因（用于日志审计）
	// required: true
	// maxLength: 200
	Reason string `json:"reason" binding:"required,max=200" example:"用户活动奖励"`
}

// AdminSetScoreResponse 管理员设置积分响应数据
// swagger:model
type AdminSetScoreResponse struct {
	UserID     uint64 `json:"user_id" example:"10086"`        // 被操作用户ID
	OldScore   int    `json:"old_score" example:"150"`        // 操作前积分
	NewScore   int    `json:"new_score" example:"200"`        // 操作后积分
	Change     int    `json:"change" example:"50"`            // 积分变化量（可为负数）
	Operation  string `json:"operation" example:"add"`        // 执行的操作类型
	OperatorID uint   `json:"operator_id" example:"1"`        // 操作管理员ID
	Reason     string `json:"reason" example:"用户活动奖励"`        // 操作原因
	Timestamp  int64  `json:"timestamp" example:"1700000000"` // Unix时间戳（秒）
}

// AdminGetUserScore 获取用户积分
// @Summary 获取用户积分
// @Description 获取指定用户积分，不传id则获取所有用户积分列表
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query int false "用户ID"
// @Success 200 {object} response.Response{data=object}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/users/score [get]
func (h *UserHandler) AdminGetUserScore(c *gin.Context) {
	// 1. 获取用户ID参数
	targetID := c.Query("id")
	if targetID == "" {
		// 如果没有id参数，返回所有用户积分列表
		users, err := h.userSvc.GetAllUsersWithScore()
		if err != nil {
			response.InternalError(c, "查询用户积分失败")
			return
		}
		response.Success(c, users)
		return
	}

	// 2. 解析用户ID
	id, err := strconv.ParseUint(targetID, 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidUserID.Error())
		return
	}

	// 3. 获取用户积分
	score, err := h.userSvc.GetScoreById(uint(id))
	if err != nil {
		response.InternalError(c, apperrors.ErrFailedToQueryScore.Error())
		return
	}

	// 4. 返回结果
	response.Success(c, gin.H{
		"user_id": id, // 建议加上 user_id，方便前端
		"score":   score,
	})
}

// 内部方法

// handleRoleError 统一处理角色变更错误，集中映射到 HTTP 响应
func handleRoleError(c *gin.Context, err error) {
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
