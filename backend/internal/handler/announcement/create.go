package announcement

import (
	"tiny-forum/internal/model/dto"
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
// @Param body body dto.CreateAnnouncementRequest true "公告信息"
// @Success 200 {object} response.Response{data=po.Announcement} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限（非管理员）"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/announcements [post]
func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req dto.CreateAnnouncementRequest
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
