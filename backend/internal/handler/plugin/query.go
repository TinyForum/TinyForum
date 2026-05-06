package plugin

import (
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/converter"
	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// List 获取插件列表
// @Summary 用户获取插件列表
// @Description 获取插件列表
// @Tags 插件管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} common.BasicResponse "获取成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "无权限"
// @Failure 500 {object} common.BasicResponse "服务器内部错误"
// @Router /plugins [get]
func (h *PluginHandler) List(c *gin.Context) {
	var req request.PluginListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Infof("绑定错误: ", err)
		response.BadRequest(c, apperrors.ErrInvalidRequest.Error())
		return
	}

	common.ApplyDefaults(&req)

	// Request -> BO
	queryBO := converter.PluginListRequestToBO(&req)

	// 调用 Service
	pageBO, err := h.service.ListPlugins(c.Request.Context(), queryBO)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// PageResult[BO] -> PageResult[VO]
	pageVO := converter.PageBOToPageVO(pageBO, converter.PluginBOToVO)

	response.SuccessPage(c, pageVO.List, pageVO.Total, pageVO.Page, pageVO.PageSize)

}
