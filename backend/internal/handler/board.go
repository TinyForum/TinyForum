package handler

import (
	"errors"
	"strconv"
	"time"

	"tiny-forum/internal/model"
	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BoardHandler struct {
	boardSvc *service.BoardService
}

func NewBoardHandler(boardSvc *service.BoardService) *BoardHandler {
	return &BoardHandler{boardSvc: boardSvc}
}

// ── Board CRUD ────────────────────────────────────────────────────────────────

// Create 创建板块
// @Summary 创建新板块
// @Description 管理员创建一个新的论坛板块，需要管理员权限
// @Tags 板块管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body service.CreateBoardInput true "板块信息"
// @Success 200 {object} response.Response{data=model.Board} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Router /boards [post]
func (h *BoardHandler) Create(c *gin.Context) {
	var input service.CreateBoardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	board, err := h.boardSvc.Create(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, board)
}

// Update 更新板块（仅管理员）
// @Summary 更新板块信息
// @Description 管理员更新指定板块的信息，需要管理员权限
// @Tags 板块管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID" minimum(1) example(1)
// @Param body body service.CreateBoardInput true "板块信息"
// @Success 200 {object} response.Response{data=model.Board} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误或板块ID无效"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/{id} [put]
func (h *BoardHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	var input service.CreateBoardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	board, err := h.boardSvc.Update(uint(id), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, board)
}

// Delete 删除板块（仅管理员）
// @Summary 删除板块（仅管理员）
// @Description 删除指定板块，需要管理员权限
// @Tags 板块管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Success 200 {object} response.Response{data=object} "删除成功"
// @Failure 400 {object} response.Response "无效的板块ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/{id} [delete]
func (h *BoardHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	if err := h.boardSvc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}

// GetByID 获取板块详情
// @Summary 获取板块详情
// @Description 根据ID获取板块详细信息
// @Tags 板块管理
// @Produce json
// @Param id path int true "板块ID"
// @Success 200 {object} response.Response{data=model.Board} "获取成功"
// @Failure 400 {object} response.Response "无效的板块ID"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/{id} [get]
func (h *BoardHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	board, err := h.boardSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, board)
}

// GetBoardBySlug 根据 slug 获取板块
// @Summary 根据板块标识符获取板块
// @Description 根据板块标识符（slug）获取板块信息
// @Tags 板块管理
// @Produce json
// @Param slug path string true "板块标识符"
// @Success 200 {object} response.Response{data=model.Board} "获取成功"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/slug/{slug} [get]
func (h *BoardHandler) GetBoardBySlug(c *gin.Context) {
	slug := c.Param("slug")
	board, err := h.boardSvc.GetBoardBySlug(slug)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, board)
}

// List 获取板块列表
// @Summary 获取板块列表
// @Description 分页获取板块列表
// @Tags 板块管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Board}} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /boards [get]
func (h *BoardHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	boards, total, err := h.boardSvc.List(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, boards, total, page, pageSize)
}

// GetTree 获取板块树形结构
// @Summary 获取板块树形结构
// @Description 获取所有板块的树形层级结构
// @Tags 板块管理
// @Produce json
// @Success 200 {object} response.Response{data=[]model.BoardTree} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /boards/tree [get]
func (h *BoardHandler) GetTree(c *gin.Context) {
	tree, err := h.boardSvc.GetTree()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, tree)
}

// GetPostsBySlug 获取板块下的帖子列表
// @Summary 获取板块下的帖子列表
// @Description 分页获取指定板块下的所有帖子
// @Tags 板块管理
// @Produce json
// @Param slug path string true "板块标识符"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Post}} "获取成功"
// @Failure 400 {object} response.Response "板块 slug 不能为空"
// @Failure 404 {object} response.Response "板块不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /boards/slug/{slug}/posts [get]
func (h *BoardHandler) GetPostsBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "板块 slug 不能为空")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	posts, total, err := h.boardSvc.GetPostsBySlug(slug, page, pageSize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "板块不存在")
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// ── 版主申请 ──────────────────────────────────────────────────────────────────

// ApplyModerator 用户申请成为版主
// @Summary 申请成为版主
// @Description 用户申请成为指定板块的版主，需要登录认证
// @Tags 版主管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param body body model.ApplyModeratorInput true "申请信息"
// @Success 200 {object} response.Response{data=object} "申请提交成功"
// @Failure 400 {object} response.Response "请求参数错误或板块ID无效"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限或已申请"
// @Router /boards/{id}/moderators/appapply-moderatorly [post]
func (h *BoardHandler) ApplyModerator(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}

	var input model.ApplyModeratorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 从 JWT 上下文注入，不信任客户端传参
	input.UserID = c.GetUint("user_id")
	input.Username = c.GetString("username")
	input.BoardID = uint(boardID)

	if err := h.boardSvc.ApplyModerator(input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "申请已提交，请等待管理员审核"})
}

// CancelApplication 用户撤销自己的版主申请
// @Summary 撤销版主申请
// @Description 用户撤销自己提交的版主申请，只能撤销状态为 pending 的申请
// @Tags 版主管理
// @Produce json
// @Security ApiKeyAuth
// @Param application_id path int true "申请ID"
// @Success 200 {object} response.Response{data=object} "撤销成功"
// @Failure 400 {object} response.Response "无效的申请ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（只能撤销自己的申请）"
// @Failure 404 {object} response.Response "申请不存在"
// @Router /boards/apply/{application_id} [delete]
func (h *BoardHandler) CancelApplication(c *gin.Context) {
	applicationID, err := strconv.ParseUint(c.Param("application_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的申请ID")
		return
	}
	userID := c.GetUint("user_id")

	if err := h.boardSvc.CancelApplication(uint(applicationID), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "申请已撤销"})
}

// GetUserApplications 获取当前用户的所有申请记录
// @Summary 获取我的申请记录
// @Description 获取当前用户提交的所有版主申请记录
// @Tags 版主管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.ModeratorApplication}} "申请列表"
// @Failure 401 {object} response.Response "未授权"
// @Router /boards/moderators/applications [get]
func (h *BoardHandler) GetUserApplications(c *gin.Context) {
	userID := c.GetUint("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	applications, total, err := h.boardSvc.GetUserApplications(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, response.PageData{
		List:     applications,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// ReviewApplicationRequest 审批请求
type ReviewApplicationRequest struct {
	Approve            bool   `json:"approve" binding:"required"`
	ReviewNote         string `json:"review_note" binding:"max=500"`
	CanDeletePost      bool   `json:"can_delete_post"`
	CanPinPost         bool   `json:"can_pin_post"`
	CanEditAnyPost     bool   `json:"can_edit_any_post"`
	CanManageModerator bool   `json:"can_manage_moderator"`
	CanBanUser         bool   `json:"can_ban_user"`
}

// ReviewApplication 管理员审批版主申请（通过或拒绝）
// @Summary 审批版主申请
// @Description 管理员审批用户的版主申请，通过时可设置版主权限，拒绝时需填写原因
// @Tags 版主管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param application_id path int true "申请ID"
// @Param body body ReviewApplicationRequest true "审批信息"
// @Success 200 {object} response.Response{data=object} "审批完成"
// @Failure 400 {object} response.Response "请求参数错误或申请ID无效"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（需要管理员权限）"
// @Failure 404 {object} response.Response "申请不存在"
// @Router /admin/boards/applications/{application_id}/review [post]
func (h *BoardHandler) ReviewApplication(c *gin.Context) {
	applicationID, err := strconv.ParseUint(c.Param("application_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的申请ID")
		return
	}

	var body struct {
		Approve            bool   `json:"approve"`
		ReviewNote         string `json:"review_note"          binding:"max=500"`
		CanDeletePost      *bool  `json:"can_delete_post"`
		CanPinPost         *bool  `json:"can_pin_post"`
		CanEditAnyPost     *bool  `json:"can_edit_any_post"`
		CanManageModerator *bool  `json:"can_manage_moderator"`
		CanBanUser         *bool  `json:"can_ban_user"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	reviewerID := c.GetUint("user_id")
	input := service.ReviewApplicationInput{
		ApplicationID:      uint(applicationID),
		Approve:            body.Approve,
		ReviewNote:         body.ReviewNote,
		CanDeletePost:      body.CanDeletePost,
		CanPinPost:         body.CanPinPost,
		CanEditAnyPost:     body.CanEditAnyPost,
		CanManageModerator: body.CanManageModerator,
		CanBanUser:         body.CanBanUser,
	}

	if err := h.boardSvc.ReviewApplication(c.Request.Context(), input, reviewerID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "审批完成"})
}

// ListApplications 管理员分页查询版主申请列表
// @Summary 获取版主申请列表
// @Description 管理员分页查询版主申请列表，可按板块和状态筛选
// @Tags 版主管理
// @Produce json
// @Security ApiKeyAuth
// @Param board_id query int false "板块ID"
// @Param status query string false "申请状态" Enums(pending, approved, rejected)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]object}} "获取成功"
// @Failure 400 {object} response.Response "无效的板块ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（需要管理员权限）"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/boards/applications [get]
func (h *BoardHandler) ListApplications(c *gin.Context) {
	var boardID *uint
	if raw := c.Query("board_id"); raw != "" {
		id, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的板块ID")
			return
		}
		uid := uint(id)
		boardID = &uid
	}

	status := model.ApplicationStatus(c.Query("status")) // 空串 = 不过滤
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	apps, total, err := h.boardSvc.ListApplications(boardID, status, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, apps, total, page, pageSize)
}

// ── 版主管理 ──────────────────────────────────────────────────────────────────

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
	input := service.AddModeratorInput{
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
	input := service.UpdateModeratorPermissionsInput{
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

// ── 禁言管理 ──────────────────────────────────────────────────────────────────

// BanUser 禁言用户
// @Summary 禁言用户
// @Description 在指定板块禁言用户，需要版主或管理员权限
// @Tags 禁言管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param body body object true "禁言信息"
// @Param body.user_id body int true "用户ID" example(10086)
// @Param body.reason body string true "禁言原因" example("发布违规内容")
// @Param body.expires_at body string false "过期时间（RFC3339格式，空表示永久）" example("2024-12-31T23:59:59Z")
// @Success 200 {object} response.Response{data=object} "禁言成功"
// @Failure 400 {object} response.Response "请求参数错误或板块ID无效"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（需要版主或管理员权限）"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/{id}/bans [post]
func (h *BoardHandler) BanUser(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}

	var body struct {
		UserID    uint   `json:"user_id"   binding:"required"`
		Reason    string `json:"reason"    binding:"required"`
		ExpiresAt string `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	banInput := service.BanUserInput{
		UserID:  body.UserID,
		BoardID: uint(boardID),
		Reason:  body.Reason,
	}
	if body.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, body.ExpiresAt)
		if err != nil {
			response.BadRequest(c, "无效的过期时间格式（需 RFC3339）")
			return
		}
		banInput.ExpiresAt = &t
	}

	operatorID := c.GetUint("user_id")
	if err := h.boardSvc.BanUser(banInput, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "禁言成功"})
}

// UnbanUser 解除禁言
// @Summary 解除禁言
// @Description 解除用户在指定板块的禁言，需要版主或管理员权限
// @Tags 禁言管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param user_id path int true "用户ID"
// @Success 200 {object} response.Response{data=object} "解除禁言成功"
// @Failure 400 {object} response.Response "无效的板块ID或用户ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（需要版主或管理员权限）"
// @Failure 404 {object} response.Response "禁言记录不存在"
// @Router /boards/{id}/bans/{user_id} [delete]
func (h *BoardHandler) UnbanUser(c *gin.Context) {
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
	if err := h.boardSvc.UnbanUser(uint(userID), uint(boardID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "解除禁言成功"})
}

// ── 帖子管理（版主） ──────────────────────────────────────────────────────────

// DeletePost 版主删除帖子
// @Summary 删除帖子（版主）
// @Description 版主或管理员删除指定板块下的帖子
// @Tags 帖子管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param post_id path int true "帖子ID"
// @Success 200 {object} response.Response{data=object} "删除成功"
// @Failure 400 {object} response.Response "无效的板块ID或帖子ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（需要版主或管理员权限）"
// @Failure 404 {object} response.Response "帖子不存在"
// @Router /boards/{id}/posts/{post_id} [delete]
func (h *BoardHandler) DeletePost(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "admin" || role == "super_admin"

	if err := h.boardSvc.DeletePost(uint(boardID), uint(postID), userID, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}

// PinPost 版主置顶/取消置顶帖子
// @Summary 置顶/取消置顶帖子
// @Description 版主或管理员置顶或取消置顶指定板块下的帖子
// @Tags 帖子管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param post_id path int true "帖子ID"
// @Param body body object true "置顶选项"
// @Param body.pin_in_board body bool true "是否置顶" example(true)
// @Success 200 {object} response.Response{data=object} "操作成功"
// @Failure 400 {object} response.Response "无效的板块ID或帖子ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（需要版主或管理员权限）"
// @Failure 404 {object} response.Response "帖子不存在"
// @Router /boards/{id}/posts/{post_id}/pin [put]
func (h *BoardHandler) PinPost(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	var body struct {
		PinInBoard bool `json:"pin_in_board"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.boardSvc.PinPost(uint(boardID), uint(postID), body.PinInBoard); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "操作成功"})
}

// ── Swagger 请求体类型声明 ────────────────────────────────────────────────────
