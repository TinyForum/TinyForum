package handler

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"tiny-forum/internal/model"
	"tiny-forum/internal/service"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	userSvc  *service.UserService
	notifSvc *service.NotificationService
}

func NewUserHandler(userSvc *service.UserService, notifSvc *service.NotificationService) *UserHandler {
	return &UserHandler{userSvc: userSvc, notifSvc: notifSvc}
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
	// 1. 参数解析
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParams(c, []response.ValidationError{
			{Field: "id", Message: "无效的用户ID格式"},
		})
		return
	}

	// 2. 获取当前登录用户（可选，未登录时为 0）
	viewerID := getViewerID(c)

	// 3. 调用服务层
	profile, err := h.userSvc.GetUserProfile(uint(targetID), viewerID)
	if err != nil {
		// 根据错误类型自动返回正确的状态码和错误信息
		// 如果是用户不存在 -> 404 + Code 20001
		// 如果是数据库错误 -> 500 + Code 50000
		response.Error(c, err)
		return
	}

	response.Success(c, profile)
}

// 辅助函数：从 context 获取当前登录用户 ID
func getViewerID(c *gin.Context) uint {
	if v, exists := c.Get("user_id"); exists {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
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

// GetCurrentUserRoleResponse 获取当前用户角色响应
type GetCurrentUserRoleResponse struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

// GetCurrentUserRole 获取当前登录用户的角色
// @Summary 获取当前用户角色
// @Description 从数据库查询当前登录用户的角色信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=GetCurrentUserRoleResponse} "操作成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "用户不存在"
// @Router /users/me/role [get]
func (h *UserHandler) GetCurrentUserRole(c *gin.Context) {
	// 1. 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未授权操作")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		response.InternalError(c, "用户身份解析失败")
		return
	}

	// 2. 从数据库查询用户角色
	role, err := h.userSvc.GetUserRoleById(userIDUint)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "用户不存在")
			return
		}
		response.InternalError(c, "查询用户角色失败")
		return
	}

	// 3. 返回响应
	response.Success(c, GetCurrentUserRoleResponse{
		UserID: userIDUint,
		Role:   role,
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

// @Summary 管理员设置用户状态
// @Tags 管理接口
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param body body object{is_active=bool} true "状态"
// @Router /admin/users/{id}/active [put]
func (h *UserHandler) AdminSetActive(c *gin.Context) {
	// 1. 解析目标用户ID
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParams(c, []response.ValidationError{
			{Field: "id", Message: "无效的用户ID格式"},
		})
		return
	}

	// 2. 获取当前操作用户ID（中间件已设置）
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

	// 3. 管理员不能修改自己的激活状态（防止误操作把自己禁用）
	if currentID == uint(targetID) {
		response.Forbidden(c, "不能修改自己的激活状态")
		return
	}

	// 4. 绑定请求体
	var body struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.InvalidParams(c, []response.ValidationError{
			{Field: "is_active", Message: "is_active 字段必须为布尔值"},
		})
		return
	}

	// 5. 调用 Service 层
	if err := h.userSvc.SetActive(uint(targetID), currentID, body.IsActive); err != nil {
		// 自动根据 err 类型返回合适的响应
		response.Error(c, err)
		return
	}

	// 6. 成功响应（建议返回更新后的用户状态，而不是简单的 message）
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
	// 1. 解析目标用户ID
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParams(c, response.SimpleValidationError("id", "无效的用户ID格式"))
		return
	}

	// 2. 获取当前操作用户ID
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

	// 3. 绑定请求体
	var body struct {
		IsBlocked bool `json:"is_blocked"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		// 打印真实错误用于调试
		fmt.Printf("❌ Bind error: %v\n", err)

		response.InvalidParams(c, response.SimpleValidationError("is_blocked", "is_blocked 字段必须为布尔值"))
		return
	}

	// 调试日志放在绑定之后
	fmt.Printf("✅ Parsed body: is_blocked=%v\n", body.IsBlocked)

	// 4. 调用 Service 层
	if err := h.userSvc.SetBlocked(uint(targetID), currentID, body.IsBlocked); err != nil {
		response.Error(c, err)
		return
	}

	// 5. 成功响应
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
	operatorID, exists := c.Get("id")
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
	// 1. 获取操作者ID
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

	// 2. 解析目标用户ID
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParams(c, response.SimpleValidationError("id", "无效的用户ID格式"))
		return
	}

	// 3. 调用 Service 删除用户
	if err := h.userSvc.DeleteUser(operatorUint, uint(targetID)); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"message":     "删除用户成功",
		"user_id":     targetID,
		"operator_id": operatorUint,
	})
}

// handler/user_handler.go

// AdminResetUserPasswordResponse 重置密码响应体
type AdminResetUserPasswordResponse struct {
	Message    string `json:"message" example:"临时密码已生成并发送给用户"`
	UserID     uint   `json:"user_id" example:"123"`
	OperatorID uint   `json:"operator_id" example:"1"`
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
	// 1. 获取操作者ID
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

	// 2. 解析目标用户ID
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.InvalidParams(c, response.SimpleValidationError("id", "无效的用户ID格式"))
		return
	}

	// 3. 调用 Service 重置密码（返回生成的临时密码）
	tempPassword, err := h.userSvc.ResetUserPasswordWithTemp(operatorUint, uint(targetID))
	if err != nil {
		response.Error(c, err)
		return
	}

	// 5. 发送通知给用户
	h.sendTempPasswordNotification(uint(targetID), operatorUint, tempPassword)

	// 6. 使用结构体返回（不返回密码，只返回成功状态）
	response.Success(c, AdminResetUserPasswordResponse{
		Message:    "临时密码已生成并发送给用户",
		UserID:     uint(targetID),
		OperatorID: operatorUint,
	})
}

// sendTempPasswordNotification 发送临时密码通知
func (h *UserHandler) sendTempPasswordNotification(targetID, operatorID uint, tempPassword string) {
	message := fmt.Sprintf(
		"管理员已重置您的密码。临时密码为：%s，有效期 30 分钟，请尽快登录并修改密码，以防被盗。",
		tempPassword,
	)

	// 发送系统通知
	h.notifSvc.Create(targetID, &operatorID, model.NotifySystem, message, nil, "")
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
