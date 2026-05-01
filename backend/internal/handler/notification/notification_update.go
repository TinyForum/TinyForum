package notification

import (
	"strconv"
	"tiny-forum/internal/model/dto"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// BatchMarkRead 批量标记通知为已读
// @Summary 批量标记通知为已读
// @Description 批量将指定ID的通知标记为已读，或不传ID则标记所有为已读
// @Tags 通知管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body dto.BatchMarkReadRequest false "批量标记请求"
// @Success 200 {object} response.Response{data=dto.BatchMarkReadResponse} "标记成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /notifications/batch/read [patch]
func (h *NotificationHandler) BatchMarkRead(c *gin.Context) {
	var req dto.BatchMarkReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 允许空body，表示标记所有
		req = dto.BatchMarkReadRequest{}
	}

	userID := c.GetUint("user_id")

	var updatedCount int64
	var err error

	if len(req.IDs) == 0 {
		// 标记所有为已读，直接返回更新的数量
		updatedCount, err = h.notifSvc.MarkAllRead(userID)
	} else {
		// 批量标记指定通知
		updatedCount, err = h.notifSvc.BatchMarkRead(userID, req.IDs)
	}

	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, dto.BatchMarkReadResponse{
		Message:      "标记成功",
		UpdatedCount: updatedCount,
	})
}

// MarkRead 标记单个通知为已读
// @Summary 标记单个通知为已读
// @Description 将指定ID的通知标记为已读
// @Tags 通知管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path uint true "通知ID"
// @Success 200 {object} response.Response{data=object} "标记成功"
// @Failure 400 {object} response.Response "无效的通知ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权操作"
// @Failure 404 {object} response.Response "通知不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /notifications/{id}/read [post]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	notifIDStr := c.Param("id")

	// 参数验证在 Handler 层
	notifID, err := strconv.ParseUint(notifIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的通知ID")
		return
	}

	// 调用 Service 层（传递正确的类型）
	if err := h.notifSvc.MarkRead(uint(notifID), userID); err != nil {
		switch err.Error() {
		case "通知不存在":
			response.NotFound(c, err.Error())
		case "无权操作此通知":
			response.Forbidden(c, err.Error())
		default:
			response.InternalError(c, err.Error())
		}
		return
	}
	response.Success(c, gin.H{"message": "已标记为已读"})
}
