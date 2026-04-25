package announcement

import (
	announcementService "tiny-forum/internal/service/announcement"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Update 更新公告
func (h *AnnouncementHandler) Update(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	var req announcementService.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.Update(c.Request.Context(), id, &req, userID); err != nil {
		handleAnnouncementServiceError(c, err)
		return
	}
	response.Success(c, nil)
}

// Publish 发布公告
func (h *AnnouncementHandler) Publish(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.Publish(c.Request.Context(), id, userID); err != nil {
		handleAnnouncementServiceError(c, err)
		return
	}
	response.Success(c, nil)
}

// Archive 归档公告
func (h *AnnouncementHandler) Archive(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.Archive(c.Request.Context(), id, userID); err != nil {
		handleAnnouncementServiceError(c, err)
		return
	}
	response.Success(c, nil)
}

// Pin 置顶/取消置顶
func (h *AnnouncementHandler) Pin(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	var req struct {
		Pinned bool `json:"pinned"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.Pin(c.Request.Context(), id, req.Pinned, userID); err != nil {
		handleAnnouncementServiceError(c, err)
		return
	}
	response.Success(c, nil)
}
