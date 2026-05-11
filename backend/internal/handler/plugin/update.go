package plugin

import (
	"strconv"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

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
