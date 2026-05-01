package announcement

import (
	"strconv"
	"tiny-forum/internal/model/po"
	"tiny-forum/internal/model/query"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetByID 获取公告详情（用户端）
// @Summary 获取公告详情
// @Description 根据ID获取已发布的公告详情（用户端调用）
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response{data=po.Announcement} "获取成功"
// @Failure 400 {object} response.Response "无效的公告ID"
// @Failure 404 {object} response.Response "公告不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /announcements/{id} [get]
func (h *AnnouncementHandler) GetByID(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}

	announcement, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, announcement)
}

// List 普通用户获取已发布的公告列表
// @Summary 获取公告列表
// @Description 分页获取已发布的公告列表（用户端）
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param board_id query int false "版块ID（可选）"
// @Success 200 {object} response.Response{data=object{list=[]po.Announcement,total=int,page=int,pageSize=int}} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /announcements [get]
func (h *AnnouncementHandler) List(c *gin.Context) {
	var req query.ListAnnouncements
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	// 强制只查已发布
	published := po.AnnouncementStatusPublished
	req.Status = &published

	resp, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, resp.Announcements, resp.Total, resp.Page, resp.PageSize)
}

// GetPinned 获取置顶公告
// @Summary 获取置顶公告
// @Description 获取置顶的公告列表，可按版块过滤
// @Tags 公告管理
// @Accept json
// @Produce json
// @Param board_id query int false "版块ID（可选）"
// @Success 200 {object} response.Response{data=[]po.Announcement} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
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
