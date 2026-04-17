package post

import (
	"errors"
	"strconv"

	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminList 管理员获取帖子列表
// @Summary 管理员获取帖子列表
// @Description 管理员分页获取所有帖子列表，支持关键词搜索
// @Tags 管理接口
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Post}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/posts [get]
func (h *PostHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	posts, total, err := h.postSvc.AdminList(page, pageSize, keyword)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// AdminTogglePin 管理员切换帖子置顶状态
// @Summary 切换帖子置顶状态
// @Description 管理员切换指定帖子的置顶状态
// @Tags 管理接口
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Success 200 {object} response.Response{data=object} "操作成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/posts/{id}/pin [put]
func (h *PostHandler) AdminTogglePin(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	if err := h.postSvc.TogglePin(uint(postID)); err != nil {
		if errors.Is(err, apperrors.ErrPostNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "操作成功"})
}
