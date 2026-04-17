// package handler

// import (
// 	"strconv"

// 	"tiny-forum/internal/service"
// 	"tiny-forum/pkg/response"

// 	"github.com/gin-gonic/gin"
// )

// type TagHandler struct {
// 	tagSvc *service.TagService
// }

// func NewTagHandler(tagSvc *service.TagService) *TagHandler {
// 	return &TagHandler{tagSvc: tagSvc}
// }

// // List 获取标签列表
// // @Summary 获取所有标签
// // @Description 获取系统中所有的标签列表
// // @Tags 标签管理
// // @Produce json
// // @Success 200 {object} response.Response{data=[]model.Tag} "获取成功"
// // @Failure 500 {object} response.Response "服务器内部错误"
// // @Router /tags [get]
// func (h *TagHandler) List(c *gin.Context) {
// 	tags, err := h.tagSvc.List()
// 	if err != nil {
// 		response.InternalError(c, err.Error())
// 		return
// 	}
// 	response.Success(c, tags)
// }

// // Create 创建标签
// // @Summary 创建标签（仅管理员）
// // @Description 创建一个新的标签，需要管理员权限
// // @Tags 标签管理
// // @Accept json
// // @Produce json
// // @Security ApiKeyAuth
// // @Param body body service.CreateTagInput true "标签信息"
// // @Success 200 {object} response.Response{data=model.Tag} "创建成功"
// // @Failure 400 {object} response.Response "请求参数错误"
// // @Failure 401 {object} response.Response "未授权"
// // @Failure 403 {object} response.Response "无权限"
// // @Failure 500 {object} response.Response "服务器内部错误"
// // @Router /tags [post]
// func (h *TagHandler) Create(c *gin.Context) {
// 	var input service.CreateTagInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		response.BadRequest(c, err.Error())
// 		return
// 	}
// 	tag, err := h.tagSvc.Create(input)
// 	if err != nil {
// 		response.InternalError(c, err.Error())
// 		return
// 	}
// 	response.Success(c, tag)
// }

// // Update 更新标签
// // @Summary 更新标签（仅管理员）
// // @Description 更新指定标签的信息，需要管理员权限
// // @Tags 标签管理
// // @Accept json
// // @Produce json
// // @Security ApiKeyAuth
// // @Param id path int true "标签ID"
// // @Param body body service.CreateTagInput true "标签信息"
// // @Success 200 {object} response.Response{data=model.Tag} "更新成功"
// // @Failure 400 {object} response.Response "请求参数错误或无效的标签ID"
// // @Failure 401 {object} response.Response "未授权"
// // @Failure 403 {object} response.Response "无权限"
// // @Failure 404 {object} response.Response "标签不存在"
// // @Failure 500 {object} response.Response "服务器内部错误"
// // @Router /tags/{id} [put]
// func (h *TagHandler) Update(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		response.BadRequest(c, "无效的标签ID")
// 		return
// 	}
// 	var input service.CreateTagInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		response.BadRequest(c, err.Error())
// 		return
// 	}
// 	tag, err := h.tagSvc.Update(uint(id), input)
// 	if err != nil {
// 		response.InternalError(c, err.Error())
// 		return
// 	}
// 	response.Success(c, tag)
// }

// // Delete 删除标签
// // @Summary 删除标签（仅管理员）
// // @Description 删除指定标签，需要管理员权限
// // @Tags 标签管理
// // @Produce json
// // @Security ApiKeyAuth
// // @Param id path int true "标签ID"
// // @Success 200 {object} response.Response{data=object} "删除成功"
// // @Failure 400 {object} response.Response "无效的标签ID"
// // @Failure 401 {object} response.Response "未授权"
// // @Failure 403 {object} response.Response "无权限"
// // @Failure 404 {object} response.Response "标签不存在"
// // @Failure 500 {object} response.Response "服务器内部错误"
// // @Router /tags/{id} [delete]
// func (h *TagHandler) Delete(c *gin.Context) {
// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		response.BadRequest(c, "无效的标签ID")
// 		return
// 	}
// 	if err := h.tagSvc.Delete(uint(id)); err != nil {
// 		response.InternalError(c, err.Error())
// 		return
// 	}
// 	response.Success(c, gin.H{"message": "删除成功"})
// }

package handler
