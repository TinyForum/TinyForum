package plugin

import (
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// TogglePlugin 切换插件启用状态
// @Summary 切换插件状态
// @Description 切换插件状态
// @Tags 插件管理
// @Security ApiKeyAuth
// @Param slug path string true "插件标识"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /plugins/{slug}/toggle [patch]
func (h *Handler) TogglePlugin(c *gin.Context) {
	// 1. 获取当前用户ID
	// userID := c.GetUint("user_id")

	// 2. 获取路径参数中的插件ID
	pluginSlug := c.Param("slug")
	var err error
	if err != nil {
		response.HandleError(c, apperrors.ErrValidation)

	}
	// 3. 调用服务层方法，更新插件状态
	err = h.svc.TogglePluginStatus(c, pluginSlug)
	if err != nil {
		response.HandleError(c, err)
	}
	response.Success(c, nil)
}
