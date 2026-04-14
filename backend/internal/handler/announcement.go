package handler

import (
	"strconv"
	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	service service.AnnouncementService
}

func NewAnnouncementHandler(service service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{service: service}
}

// Create 创建公告
// @Summary 创建公告
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param request body service.CreateAnnouncementRequest true "创建公告请求"
// @Success 200 {object} response.Response{data=model.Announcement}
// @Router /admin/announcements [post]
func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req service.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id") // 从上下文中获取用户ID
	announcement, err := h.service.Create(c.Request.Context(), &req, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, announcement)
}

// Update 更新公告
// @Summary 更新公告
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Param request body service.UpdateAnnouncementRequest true "更新公告请求"
// @Success 200 {object} response.Response
// @Router /admin/announcements/{id} [put]
func (h *AnnouncementHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的公告ID")
		return
	}

	var req service.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.Update(c.Request.Context(), uint(id), &req, userID); err != nil {
		if err == service.ErrAnnouncementNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除公告
// @Summary 删除公告
// @Tags 公告管理
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response
// @Router /admin/announcements/{id} [delete]
func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的公告ID")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.Delete(c.Request.Context(), uint(id), userID); err != nil {
		if err == service.ErrAnnouncementNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID 获取公告详情
// @Summary 获取公告详情
// @Tags 公告管理
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response{data=model.Announcement}
// @Router /announcements/{id} [get]
func (h *AnnouncementHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的公告ID")
		return
	}

	announcement, err := h.service.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == service.ErrAnnouncementNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, announcement)
}

// List 获取公告列表
// @Summary 获取公告列表
// @Tags 公告管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param board_id query int false "板块ID"
// @Param type query string false "公告类型" Enums(normal,important,emergency,event)
// @Param status query string false "状态" Enums(draft,published,archived)
// @Param is_pinned query bool false "是否置顶"
// @Param is_global query bool false "是否全局"
// @Param keyword query string false "关键词"
// @Success 200 {object} response.Response{data=service.ListAnnouncementResponse}
// @Router /announcements [get]
func (h *AnnouncementHandler) List(c *gin.Context) {
	var req service.ListAnnouncementRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// 使用分页响应格式
	response.SuccessPage(c, resp.Announcements, resp.Total, resp.Page, resp.PageSize)
}

// GetPinned 获取置顶公告
// @Summary 获取置顶公告
// @Tags 公告管理
// @Produce json
// @Param board_id query int false "板块ID"
// @Success 200 {object} response.Response{data=[]model.Announcement}
// @Router /announcements/pinned [get]
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

// Publish 发布公告
// @Summary 发布公告
// @Tags 公告管理
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response
// @Router /admin/announcements/{id}/publish [post]
func (h *AnnouncementHandler) Publish(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的公告ID")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.Publish(c.Request.Context(), uint(id), userID); err != nil {
		if err == service.ErrAnnouncementNotFound {
			response.NotFound(c, err.Error())
			return
		}
		if err == service.ErrInvalidPublishTime {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Archive 归档公告
// @Summary 归档公告
// @Tags 公告管理
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response
// @Router /admin/announcements/{id}/archive [post]
func (h *AnnouncementHandler) Archive(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的公告ID")
		return
	}

	userID := c.GetUint("user_id")
	if err := h.service.Archive(c.Request.Context(), uint(id), userID); err != nil {
		if err == service.ErrAnnouncementNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// Pin 置顶/取消置顶公告
// @Summary 置顶/取消置顶公告
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Param request body object true "置顶状态" example({"pinned":true})
// @Success 200 {object} response.Response
// @Router /admin/announcements/{id}/pin [put]
func (h *AnnouncementHandler) Pin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的公告ID")
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
	if err := h.service.Pin(c.Request.Context(), uint(id), req.Pinned, userID); err != nil {
		if err == service.ErrAnnouncementNotFound {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
