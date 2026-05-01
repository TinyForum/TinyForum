package announcement

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

// AdminList 管理员获取所有状态的公告列表

// List 列出公告
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
//
// Deprecated: 迁移到 adminHandler.ListAnnouncements
func (h *AnnouncementHandler) AdminList(c *gin.Context) {
	var req request.ListAnnouncements
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	allStatus := do.AnnouncementStatus(do.AnnouncementStatusAll)
	req.Status = &allStatus

	resp, err := h.service.List(c.Request.Context(), &req)
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
