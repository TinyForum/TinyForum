package admin

import (
	"strconv"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

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
// @Success 200 {object} common.BasicResponse  "获取成功"
// @Failure 400 {object} common.BasicResponse"无效的板块ID"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（需要管理员权限）"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /admin/boards/applications [get]
// Desp
func (h *AdminHandler) ListApplications(c *gin.Context) {
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

	status := do.ApplicationStatus(c.Query("status")) // 空串 = 不过滤
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	apps, total, err := h.service.ListApplications(boardID, status, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, apps, total, page, pageSize)
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
// @Success 200 {object} common.BasicResponse  "审批完成"
// @Failure 400 {object} common.BasicResponse"请求参数错误或申请ID无效"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（需要管理员权限）"
// @Failure 404 {object} common.BasicResponse"申请不存在"
// @Router /admin/boards/applications/{application_id}/review [post]
func (h *AdminHandler) ReviewApplication(c *gin.Context) {
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
	input := request.ReviewApplicationRequest{
		ApplicationID:      uint(applicationID),
		Approve:            body.Approve,
		ReviewNote:         body.ReviewNote,
		CanDeletePost:      body.CanDeletePost,
		CanPinPost:         body.CanPinPost,
		CanEditAnyPost:     body.CanEditAnyPost,
		CanManageModerator: body.CanManageModerator,
		CanBanUser:         body.CanBanUser,
	}

	if err := h.service.ReviewApplication(c.Request.Context(), input, reviewerID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "审批完成"})
}
