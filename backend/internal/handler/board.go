package handler

import (
	"strconv"
	"time"

	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type BoardHandler struct {
	boardSvc *service.BoardService
}

func NewBoardHandler(boardSvc *service.BoardService) *BoardHandler {
	return &BoardHandler{boardSvc: boardSvc}
}

// Create 创建板块
// @Summary 创建板块（仅管理员）
// @Description 创建一个新的板块，需要管理员权限
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

// Update 更新板块
// @Summary 更新板块（仅管理员）
// @Description 更新指定板块的信息，需要管理员权限
// @Tags 板块管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param body body service.CreateBoardInput true "板块信息"
// @Success 200 {object} response.Response{data=model.Board} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
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

// Delete 删除板块
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

// GetBySlug 根据Slug获取板块
// @Summary 根据Slug获取板块
// @Description 根据板块标识符（slug）获取板块信息
// @Tags 板块管理
// @Produce json
// @Param slug path string true "板块标识符"
// @Success 200 {object} response.Response{data=model.Board} "获取成功"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/slug/{slug} [get]
func (h *BoardHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	board, err := h.boardSvc.GetBySlug(slug)
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

// GetTree 获取板块树
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

// GetPosts 获取板块下的帖子
// @Summary 获取板块下的帖子列表
// @Description 分页获取指定板块下的所有帖子
// @Tags 板块管理
// @Produce json
// @Param id path int true "板块ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Post}} "获取成功"
// @Failure 400 {object} response.Response "无效的板块ID"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /boards/{id}/posts [get]
func (h *BoardHandler) GetPosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	posts, total, err := h.boardSvc.GetPosts(uint(id), page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// AddModerator 添加版主
// @Summary 添加版主
// @Description 为指定板块添加版主，需要版主管理权限
// @Tags 版主管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param body body AddModeratorRequest true "版主信息"
// @Success 200 {object} response.Response{data=object} "添加成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Router /boards/{id}/moderators [post]
func (h *BoardHandler) AddModerator(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}

	var input struct {
		UserID             uint `json:"user_id" binding:"required"`
		CanDeletePost      bool `json:"can_delete_post"`
		CanPinPost         bool `json:"can_pin_post"`
		CanEditAnyPost     bool `json:"can_edit_any_post"`
		CanManageModerator bool `json:"can_manage_moderator"`
		CanBanUser         bool `json:"can_ban_user"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	modInput := service.AddModeratorInput{
		UserID:             input.UserID,
		BoardID:            uint(boardID),
		CanDeletePost:      input.CanDeletePost,
		CanPinPost:         input.CanPinPost,
		CanEditAnyPost:     input.CanEditAnyPost,
		CanManageModerator: input.CanManageModerator,
		CanBanUser:         input.CanBanUser,
	}

	if err := h.boardSvc.AddModerator(modInput, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "添加版主成功"})
}

// RemoveModerator 移除版主
// @Summary 移除版主
// @Description 移除指定板块的版主，需要版主管理权限
// @Tags 版主管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param user_id path int true "用户ID"
// @Success 200 {object} response.Response{data=object} "移除成功"
// @Failure 400 {object} response.Response "无效的ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
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

	if err := h.boardSvc.RemoveModerator(uint(userID), uint(boardID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "移除版主成功"})
}

// GetModerators 获取版主列表
// @Summary 获取板块版主列表
// @Description 获取指定板块的所有版主信息
// @Tags 版主管理
// @Produce json
// @Security ApiKeyAuth
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

// BanUser 禁言用户
// @Summary 禁言用户
// @Description 在指定板块禁言用户，需要版主权限
// @Tags 禁言管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param body body BanUserRequest true "禁言信息"
// @Success 200 {object} response.Response{data=object} "禁言成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Router /boards/{id}/bans [post]
func (h *BoardHandler) BanUser(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}

	var input struct {
		UserID    uint   `json:"user_id" binding:"required"`
		Reason    string `json:"reason" binding:"required"`
		ExpiresAt string `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operatorID := c.GetUint("user_id")
	banInput := service.BanUserInput{
		UserID:  input.UserID,
		BoardID: uint(boardID),
		Reason:  input.Reason,
	}

	if input.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339, input.ExpiresAt)
		if err != nil {
			response.BadRequest(c, "无效的过期时间格式")
			return
		}
		banInput.ExpiresAt = &expiresAt
	}

	if err := h.boardSvc.BanUser(banInput, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "禁言成功"})
}

// UnbanUser 解除禁言
// @Summary 解除禁言
// @Description 解除用户在指定板块的禁言，需要版主权限
// @Tags 禁言管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param user_id path int true "用户ID"
// @Success 200 {object} response.Response{data=object} "解除成功"
// @Failure 400 {object} response.Response "无效的ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
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

// AddModeratorRequest 添加版主请求参数
type AddModeratorRequest struct {
	UserID             uint `json:"user_id" example:"1" binding:"required"`
	CanDeletePost      bool `json:"can_delete_post" example:"true"`
	CanPinPost         bool `json:"can_pin_post" example:"true"`
	CanEditAnyPost     bool `json:"can_edit_any_post" example:"false"`
	CanManageModerator bool `json:"can_manage_moderator" example:"false"`
	CanBanUser         bool `json:"can_ban_user" example:"true"`
}

// BanUserRequest 禁言用户请求参数
type BanUserRequest struct {
	UserID    uint   `json:"user_id" example:"1" binding:"required"`
	Reason    string `json:"reason" example:"发布违规内容" binding:"required"`
	ExpiresAt string `json:"expires_at" example:"2024-12-31T23:59:59Z"`
}

// DeletePost 删除板块中的帖子（版主）
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
	isAdmin := role == "admin"

	if err := h.boardSvc.DeletePost(uint(boardID), uint(postID), userID, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}

// PinPost 置顶/取消置顶帖子（版主）
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

	var input struct {
		PinInBoard bool `json:"pin_in_board"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.boardSvc.PinPost(uint(boardID), uint(postID), input.PinInBoard); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "操作成功"})
}
