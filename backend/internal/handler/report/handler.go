package report

import (
	"fmt"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/service/reports"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ReportHandler 举报处理器
type ReportHandler struct {
	service reports.ReportsService
}

// NewReportHandler 创建举报处理器
func NewReportHandler(service reports.ReportsService) *ReportHandler {
	return &ReportHandler{service: service}
}

// CreateReport 提交举报
// @Summary 提交举报
// @Description 对帖子/评论/用户提交举报
// @Tags 举报
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.CreateReportRequest true "举报内容"
// @Success 201 {object} common.BasicResponse "提交成功"
// @Failure 400 {object} common.BasicResponse "请求参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 500 {object} common.BasicResponse "服务器内部错误"
// @Router /reports [post]
func (h *ReportHandler) CreateReport(c *gin.Context) {
	var req request.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	reporterID := c.GetUint("user_id")

	// 举报人IP由服务端取真实客户端 IP，不信任客户端传入
	req.ReporterIP = c.ClientIP()

	report, err := h.service.CreateReport(c.Request.Context(), reporterID, &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	location := fmt.Sprintf("/api/v1/reports/%d", report.ID)
	response.Created(c, report, location)
}
