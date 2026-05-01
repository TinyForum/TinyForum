package admin

import (
	"errors"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ListAnnouncements 列出公告
// @Summary 列出公告
// @Description 管理员列出所有公告，需要认证及管理员权限
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.ListAnnouncements true "公告信息"
// @Success 200 {object} response.Response{data=do.Announcement} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（非管理员）"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/announcements/list [get]
func (h *AdminHandler) ListAnnouncements(c *gin.Context) {
	var req request.ListAnnouncements
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	allStatus := do.AnnouncementStatus(do.AnnouncementStatusAll)
	req.Status = &allStatus

	resp, err := h.service.ListAnnouncements(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrAnnouncementNotFound):
			response.HandleError(c, apperrors.ErrAnnouncementNotFound)
		case errors.Is(err, apperrors.ErrInsufficientPermission):
			response.HandleError(c, apperrors.ErrInsufficientPermission)
		default:
			// 记录未知错误日志
			logger.Error("unexpected error", zap.Error(err))
			response.InternalError(c, apperrors.ErrSystemBusy.Message)
		}
		return
	}
	response.SuccessPage(c, resp.Announcements, resp.Total, resp.Page, resp.PageSize)
}

// Create 创建公告
// @Summary 创建公告
// @Description 管理员创建一条新公告，需要认证及管理员权限
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.CreateAnnouncement true "公告信息"
// @Success 200 {object} response.Response{data=do.Announcement} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（非管理员）"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/announcements [post]
func (h *AdminHandler) CreateAnnouncement(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.CreateAnnouncement
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorf("绑定请求失败: user: %d result error: %v",userID,err)
		response.BadRequest(c, err.Error())
		return
	}
	logger.Infof("用户请求创建公告: user: %d request: %v",userID,req)
	

	
		
	announcement, err := h.service.CreateAnnouncement(c.Request.Context(), &req, userID)
	
	if err != nil {
		logger.Errorf("创建公告失败: user: %d result: %v",userID,err)
		response.InternalError(c, err.Error())
		return
	}
	logger.Infof("创建公告: user: %d result: %v",userID,announcement.ID)
	response.Success(c, announcement)
}

// @Summary 更新公告
// @Description 更新公告信息
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "公告ID"
// @Param body body dto.UpdateAnnouncementRequest true "公告信息"
// @Success 200 {object} response.Response "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "公告不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/announcements/{id} [put]
func (h *AdminHandler) UpdateAnnouncement(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	var req request.UpdateAnnouncement
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.UpdateAnnouncement(c.Request.Context(), id, &req, userID); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// Delete 删除公告
// @Summary 删除公告
// @Description 管理员根据ID删除公告，需要认证及管理员权限
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "公告ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 400 {object} response.Response "参数错误（无效的公告ID）"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（非管理员）"
// @Failure 404 {object} response.Response "公告不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/announcements/{id} [delete]
func (h *AdminHandler) DeleteAnnouncement(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.DeleteAnnouncement(c.Request.Context(), id, userID); err != nil {
		response.HandleError(c, err)
		logger.Error("delete announcement failed",
			zap.Uint("id", id),
			zap.Uint("user_id", userID),
			zap.Error(err),
		)
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
func (h *AdminHandler) PublishAnnouncement(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.PublishAnnouncement(c.Request.Context(), id, userID); err != nil {
		response.HandleError(c, err)
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
func (h *AdminHandler) ArchiveAnnouncement(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.ArchiveAnnouncement(c.Request.Context(), id, userID); err != nil {
		response.HandleError(c, err)
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
func (h *AdminHandler) PinAnnouncement(c *gin.Context) {
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
	if err := h.service.PinAnnouncement(c.Request.Context(), id, req.Pinned, userID); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}
