package tag

import (
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// List 获取标签列表
// @Summary 获取所有标签
// @Description 获取系统中所有的标签列表
// @Tags 标签管理
// @Produce json
// @Success 200 {object} response.Response{data=[]model.Tag} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags [get]
func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.tagSvc.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, tags)
}
