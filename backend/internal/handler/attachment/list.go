package attachment

import (
	"tiny-forum/internal/model/converter"
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetPluginFiles 我的插件
// @Summary 我的插件
// @Tags 上传管理
// @Produce json
// @Success 200 {object} common.BasicResponse
// @Router /attachment/plugin [get]
func (h *UploadHandler) ListMyPlugins(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
	}
	var request request.PluginListRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Infof("绑定错误: ", err)
		response.BadRequest(c, "参数错误")
		return
	}

	// var req request.PluginListRequest
	// if err := c.ShouldBindQuery(&req); err != nil {
	// 	logger.Infof("绑定错误: ", err)
	// 	response.BadRequest(c, apperrors.ErrInvalidRequest.Error())
	// 	return
	// }
	queryBO := converter.PluginListRequestToUserPluginBO(&request, userID.(uint))

	// resultBO, err := h.service.ListUserPlugins(c, *query)
	// if err != nil {
	// 	response.HandleError(c, err)
	// 	return
	// }

	// Request -> BO
	// queryBO := converter.PluginListRequestToBO(&req)

	// 调用 Service
	pageVO, err := h.service.ListUserPlugins(c.Request.Context(), queryBO)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	// PageResult[BO] -> PageResult[VO]
	// pageVO := converter.PageBOToPageVO(pageBO, converter.PluginBOToVO)

	response.SuccessPage(c, pageVO.List, pageVO.Total, pageVO.Page, pageVO.PageSize)
	// response.SuccessPage(c, resultBO.List, resultBO.Total, request.Page, request.PageSize)
}
