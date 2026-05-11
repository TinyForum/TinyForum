package plugin

import (
	"strconv"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary 切换插件状态
// @Description 切换插件状态
// @Tags 插件管理
// @Param id path int true "插件ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /plugin/{id}/toggle [put]

func (h *Handler) TogglePlugin(c *gin.Context) {
	// 1. 获取当前用户ID
	// userID := c.GetUint("user_id")

	// 2. 获取路径参数中的插件ID
	pluginIDStr := c.Param("id")
	pluginID, err := strconv.ParseUint(pluginIDStr, 10, 32)
	if err != nil {
		response.HandleError(c, apperrors.ErrValidation)
	}
	// 3. 调用服务层方法，更新插件状态
	err = h.svc.TogglePluginStatus(c, uint(pluginID))
	if err != nil {
		response.HandleError(c, err)
	}
	response.Success(c, nil)
}
