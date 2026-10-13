package reviewer

import (
	"strconv"

	"tiny-forum/internal/model/request"
	reviewerService "tiny-forum/internal/service/reviewer"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ReviewerHandler 审核员处理器
type ReviewerHandler struct {
	service reviewerService.ReviewerService
}

// NewReviewerHandler 创建审核员处理器
func NewReviewerHandler(service reviewerService.ReviewerService) *ReviewerHandler {
	return &ReviewerHandler{service: service}
}

// ListPendingTasks 获取待审核任务列表
// @Summary 获取待审核任务列表
// @Tags 审核
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} common.BasicResponse
// @Router /reviewer/audit/tasks [get]
func (h *ReviewerHandler) ListPendingTasks(c *gin.Context) {
	page, pageSize := parsePage(c)
	result, err := h.service.ListPendingTasks(c.Request.Context(), page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// ListReviewedTasks 获取审核历史（已处理任务）
// @Summary 获取审核历史
// @Tags 审核
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} common.BasicResponse
// @Router /reviewer/audit/history [get]
func (h *ReviewerHandler) ListReviewedTasks(c *gin.Context) {
	page, pageSize := parsePage(c)
	result, err := h.service.ListReviewedTasks(c.Request.Context(), page, pageSize)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// parsePage 解析分页参数，修正越界值
func parsePage(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}

// ApproveTask 审核通过
// @Summary 审核通过
// @Tags 审核
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "任务ID"
// @Success 200 {object} common.BasicResponse
// @Router /reviewer/audit/tasks/{id}/approve [put]
func (h *ReviewerHandler) ApproveTask(c *gin.Context) {
	h.resolve(c, true)
}

// RejectTask 审核拒绝
// @Summary 审核拒绝
// @Tags 审核
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "任务ID"
// @Success 200 {object} common.BasicResponse
// @Router /reviewer/audit/tasks/{id}/reject [put]
func (h *ReviewerHandler) RejectTask(c *gin.Context) {
	h.resolve(c, false)
}

// resolve 统一处理通过/拒绝逻辑
func (h *ReviewerHandler) resolve(c *gin.Context, approved bool) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}

	var req request.ReviewTaskRequest
	_ = c.ShouldBindJSON(&req)

	reviewerID := c.GetUint("user_id")
	if err := h.service.ReviewTask(c.Request.Context(), uint(taskID), reviewerID, approved, req.Note); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"reviewed": true})
}
