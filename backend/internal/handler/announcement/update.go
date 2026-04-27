package announcement

import (
	announcementService "tiny-forum/internal/service/announcement"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary 更新公告
// @Description 更新公告信息
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "公告ID"
// @Param body body announcementService.UpdateAnnouncementRequest true "公告信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "公告不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/announcements/{id} [put]
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

// @Summary 发布公告
// @Description 将公告状态设为已发布
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response "发布成功"
// @Failure 400 {object} response.Response "参数错误或发布时间无效"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "公告不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/announcements/{id}/publish [post]
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

// @Summary 归档公告
// @Description 将公告状态设为已归档
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response "归档成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "公告不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/announcements/{id}/archive [post]
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

// @Summary 置顶/取消置顶公告
// @Description 设置公告的置顶状态
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "公告ID"
// @Param body body object true "置顶状态" example({"pinned": true})
// @Success 200 {object} response.Response "操作成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "公告不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/announcements/{id}/pin [put]
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
