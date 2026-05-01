package board

import (
	"strconv"

	"tiny-forum/internal/model/po"
	boardService "tiny-forum/internal/service/board"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

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
// @Router /boards/{id}/moderators/apply [post]
func (h *BoardHandler) ApplyModerator(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}

	var input po.ApplyModeratorInput
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
// @Success 200 {object} response.Response{data=response.PageData{list=[]po.ModeratorApplication}} "申请列表"
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
	input := boardService.ReviewApplicationInput{
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

	status := po.ApplicationStatus(c.Query("status")) // 空串 = 不过滤
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	apps, total, err := h.boardSvc.ListApplications(boardID, status, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, apps, total, page, pageSize)
}
