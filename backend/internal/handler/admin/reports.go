package admin

import (
	"strconv"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListReports 列出所有举报内容
// @Summary 列出所有举报内容
// @Description 管理员查看所有的举报信息
// @Tags 举报管理
// @Produce json
// @Security ApiKeyAuth
// @Param status query string false "举报状态" Enums(pending, resolved, rejected)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "关键词"
// @Success 200 {object} common.BasicResponse  "获取成功"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（需要管理员权限）"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /admin/reports [get]
func (h *AdminHandler) ListReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := do.ReportStatus(c.Query("status")) // 空串 = 不过滤
	keyword := c.Query("keyword")

	query := &common.PageQuery[bo.ListReportBO]{
		Page:     page,
		PageSize: pageSize,
		Keyword:  keyword,
		Data: bo.ListReportBO{
			Status: status,
		},
	}

	reports, total, err := h.service.ListReports(c.Request.Context(), query)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.SuccessPage(c, reports, total, page, pageSize)
}
