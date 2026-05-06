package announcement

import (
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Create 创建公告
// @Summary 创建公告
// @Description 管理员创建一条新公告，需要认证及管理员权限
// @Tags 公告管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.CreateAnnouncement true "公告信息"
// @Success 200 {object} common.BasicResponse "创建成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（非管理员）"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /admin/announcements [post]
//
// Deprecated: 迁移到 adminHandler.CreateAnnouncements
func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req request.CreateAnnouncement
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
