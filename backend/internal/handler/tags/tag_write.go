package tag

import (
	"strconv"

	tagService "tiny-forum/internal/service/tag"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Create 创建标签
// @Summary 创建标签（仅管理员）
// @Description 创建一个新的标签，需要管理员权限
// @Tags 标签管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body tag.CreateTagInput true "标签信息"
// @Success 200 {object} common.BasicResponse "创建成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /tags [post]
func (h *TagHandler) Create(c *gin.Context) {
	var input tagService.CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tag, err := h.tagSvc.Create(input)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, tag)
}

// Update 更新标签
// @Summary 更新标签（仅管理员）
// @Description 更新指定标签的信息，需要管理员权限
// @Tags 标签管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "标签ID"
// @Param body body tag.CreateTagInput true "标签信息"
// @Success 200 {object} common.BasicResponse "更新成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误或无效的标签ID"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限"
// @Failure 404 {object} common.BasicResponse"标签不存在"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /tags/{id} [put]
func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}
	var input tagService.CreateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tag, err := h.tagSvc.Update(uint(id), input)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, tag)
}

// Delete 删除标签
// @Summary 删除标签（仅管理员）
// @Description 删除指定标签，需要管理员权限
// @Tags 标签管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "标签ID"
// @Success 200 {object} common.BasicResponse  "删除成功"
// @Failure 400 {object} common.BasicResponse"无效的标签ID"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限"
// @Failure 404 {object} common.BasicResponse"标签不存在"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /tags/{id} [delete]
func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}
	if err := h.tagSvc.Delete(uint(id)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}
