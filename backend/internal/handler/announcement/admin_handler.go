package announcement

import (
	"tiny-forum/internal/model"
	announcementService "tiny-forum/internal/service/announcement"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Create 创建公告
func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req announcementService.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("user_id")
	announcement, err := h.service.Create(c.Request.Context(), &req, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, announcement)
}

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

// Delete 删除公告
func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
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

// AdminList 管理员获取所有状态的公告列表
func (h *AnnouncementHandler) AdminList(c *gin.Context) {
	var req announcementService.ListAnnouncementRequest
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
