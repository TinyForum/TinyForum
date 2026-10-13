package plugin

import (
	"tiny-forum/internal/model/request"
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

	// 2. 获取路由参数中的插件ID
	pluginSlug := c.Param("slug")
	// 3. 调用服务层方法，更新插件状态
	err := h.svc.TogglePluginStatus(c, pluginSlug)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// GetPlugin 获取单个插件详情
// @Summary 获取插件详情
// @Description 根据 slug 获取插件详细信息
// @Tags 插件管理
// @Produce json
// @Security ApiKeyAuth
// @Param slug path string true "插件标识"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /plugins/{slug} [get]
func (h *Handler) GetPlugin(c *gin.Context) {
	pluginSlug := c.Param("slug")
	meta, err := h.svc.GetPlugin(c.Request.Context(), pluginSlug)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, meta)
}

// UpdatePlugin 更新插件元数据（管理员）
// @Summary 更新插件元数据
// @Description 管理员更新插件元数据（名称、版本、描述、启用状态等）
// @Tags 插件管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param slug path string true "插件标识"
// @Param body body request.UpdatePluginRequest true "更新请求"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /plugins/{slug} [put]
func (h *Handler) UpdatePlugin(c *gin.Context) {
	pluginSlug := c.Param("slug")
	var req request.UpdatePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, apperrors.ErrValidation)
		return
	}
	if err := h.svc.UpdatePlugin(c.Request.Context(), pluginSlug, &req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}
