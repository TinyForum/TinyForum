package announcement

import (
	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Delete 删除公告
// @Summary 删除公告
// @Description 管理员根据ID删除公告，需要认证及管理员权限
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "公告ID"
// @Success 200 {object} common.BasicResponse"删除成功"
// @Failure 400 {object} common.BasicResponse"参数错误（无效的公告ID）"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（非管理员）"
// @Failure 404 {object} common.BasicResponse"公告不存在"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /admin/announcements/{id} [delete]
//
// Deprecated: 迁移到 adminHandler.DeleteAnnouncements
func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id, ok := parseAnnouncementID(c)
	if !ok {
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
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
