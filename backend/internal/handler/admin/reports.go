package admin

import (
	"strconv"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary 列出所有举报内容
// @Description 管理员查看所有的举报信息
// @Tags 举报管理
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
func (h *AdminHandler) ListReports(c *gin.Context) {
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
