package announcement

import (
	"strconv"
	"tiny-forum/internal/model"
	announcementService "tiny-forum/internal/service/announcement"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetByID 获取公告详情（用户端）
func (h *AnnouncementHandler) GetByID(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}

	announcement, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		handleAnnouncementServiceError(c, err)
		return
	}
	response.Success(c, announcement)
}

// List 普通用户获取已发布的公告列表
func (h *AnnouncementHandler) List(c *gin.Context) {
	var req announcementService.ListAnnouncementRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	// 强制只查已发布
	published := model.AnnouncementStatusPublished
	req.Status = &published

	resp, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, resp.Announcements, resp.Total, resp.Page, resp.PageSize)
}

// GetPinned 获取置顶公告
func (h *AnnouncementHandler) GetPinned(c *gin.Context) {
	var boardID *uint
	if idStr := c.Query("board_id"); idStr != "" {
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err == nil {
			boardID = new(uint)
			*boardID = uint(id)
		}
	}
	announcements, err := h.service.GetPinned(c.Request.Context(), boardID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, announcements)
}
