package announcement

import (
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminList 管理员获取所有状态的公告列表
func (h *AnnouncementHandler) AdminList(c *gin.Context) {
	var req dto.ListAnnouncementRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	// 管理员：查询所有公告（状态为 "all"）
	allStatus := model.AnnouncementStatus("all")
	req.Status = &allStatus

	resp, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, resp.Announcements, resp.Total, resp.Page, resp.PageSize)
}
