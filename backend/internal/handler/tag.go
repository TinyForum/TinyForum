package handler

import (
	"strconv"

	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagSvc *service.TagService
}

func NewTagHandler(tagSvc *service.TagService) *TagHandler {
	return &TagHandler{tagSvc: tagSvc}
}

func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.tagSvc.List()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, tags)
}

func (h *TagHandler) Create(c *gin.Context) {
	var input service.CreateTagInput
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

func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的标签ID")
		return
	}
	var input service.CreateTagInput
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
