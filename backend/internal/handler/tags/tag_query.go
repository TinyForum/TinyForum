package tag

import (
	"strconv"
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

// Get 获取单个标签
// @Summary 获取单个标签
// @Description 根据标签ID获取标签信息
// @Tags 标签管理
// @Param id path int true "标签ID"
// @Produce json
// @Success 200 {object} response.Response{data=model.Tag} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /tags/{id} [get]
func (h *TagHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}
	tag, err := h.tagSvc.Get(uint(id))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, tag)
}
