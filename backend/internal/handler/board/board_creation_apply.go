package board

import (
	"strconv"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ApplyBoardCreation 用户申请创建新板块
// @Summary 申请创建板块
// @Description 用户提交创建新板块的申请，需要登录认证，由管理员审批
// @Tags 板块创建申请
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.ApplyBoardCreationRequest true "申请信息"
// @Success 200 {object} common.BasicResponse "申请提交成功"
// @Failure 400 {object} common.BasicResponse "请求参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /boards/creation-applications [post]
func (h *BoardHandler) ApplyBoardCreation(c *gin.Context) {
	var req request.ApplyBoardCreationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	uid, err := getAuthenticatedUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	req.UserID = uid

	if err := h.boardSvc.ApplyBoardCreation(c.Request.Context(), req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, common.ResponseMessage{Message: "申请已提交，请等待管理员审核"})
}

// GetUserBoardApplications 获取当前用户的板块创建申请记录
// @Summary 获取我的板块创建申请
// @Description 获取当前用户提交的所有板块创建申请记录
// @Tags 板块创建申请
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} common.BasicResponse "申请列表"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /boards/creation-applications/my [get]
func (h *BoardHandler) GetUserBoardApplications(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	apps, total, err := h.boardSvc.GetUserBoardApplications(userID, page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, common.PageResult[do.BoardCreationApplication]{
		List:     apps,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// ListBoardApplications 管理员分页查询板块创建申请列表
// @Summary 获取板块创建申请列表
// @Description 管理员分页查询板块创建申请列表，可按状态筛选
// @Tags 板块创建申请
// @Produce json
// @Security ApiKeyAuth
// @Param status query string false "申请状态" Enums(pending, approved, rejected)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} common.BasicResponse "获取成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "无权限（需要管理员权限）"
// @Router /boards/creation-applications [get]
func (h *BoardHandler) ListBoardApplications(c *gin.Context) {
	status := do.ApplicationStatus(c.Query("status")) // 空串 = 不过滤
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	apps, total, err := h.boardSvc.ListBoardApplications(status, page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.SuccessPage(c, apps, total, page, pageSize)
}

// ReviewBoardCreationApplication 管理员审批板块创建申请
// @Summary 审批板块创建申请
// @Description 管理员审批板块创建申请，通过后自动创建对应板块
// @Tags 板块创建申请
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "申请ID"
// @Param body body request.ReviewBoardCreationRequest true "审批信息"
// @Success 200 {object} common.BasicResponse "审批完成"
// @Failure 400 {object} common.BasicResponse "请求参数错误或申请ID无效"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "无权限（需要管理员权限）"
// @Failure 404 {object} common.BasicResponse "申请不存在"
// @Router /boards/creation-applications/{id}/review [post]
func (h *BoardHandler) ReviewBoardCreationApplication(c *gin.Context) {
	applicationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || applicationID == 0 {
		response.BadRequest(c, "无效的申请ID")
		return
	}

	var body struct {
		Approve    bool   `json:"approve"`
		ReviewNote string `json:"review_note" binding:"max=500"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.HandleError(c, err)
		return
	}

	reviewerID := c.GetUint("user_id")
	input := request.ReviewBoardCreationRequest{
		ApplicationID: uint(applicationID),
		Approve:       body.Approve,
		ReviewNote:    body.ReviewNote,
	}
	if err := h.boardSvc.ReviewBoardCreationApplication(c.Request.Context(), input, reviewerID); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, common.ResponseMessage{Message: "审批完成"})
}
