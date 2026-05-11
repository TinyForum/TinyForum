package admin

import (
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminList 管理员获取帖子列表
// @Summary 管理员获取帖子列表
// @Description 管理员分页获取所有帖子列表，支持关键词搜索
// @Tags 管理员后台
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} common.BasicResponse  "获取成功"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限"
// @Failure 500 {object} common.BasicResponse"服务器内部错误"
// @Router /admin/posts [get]
//
// Deprecated: 迁移到 adminHandler.ListPosts
func (h *AdminHandler) ListPosts(c *gin.Context) {
	// page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	// pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	// keyword := c.Query("keyword")
	var req request.ListPosts
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "invalid parameters: "+err.Error())
		return
	}

	postStatus := do.ParsePostStatus(req.PostStatus)
	postType := do.ParsePostType(req.PostType)
	moderationStatus := do.ParseModerationStatus(req.ModerationStatus)

	listPostsBO := &common.PageQuery[bo.ListPosts]{
		Page:     req.Page,
		PageSize: req.PageSize,
		Data: bo.ListPosts{
			PostStatus:       postStatus,
			Keyword:          req.Keyword,
			Type:             postType,
			SortBy:           req.SortBy,
			AuthorID:         req.AuthorID,
			TagNames:         req.TagNames,
			ModerationStatus: moderationStatus,
		},
	}
	posts, total, err := h.service.ListPosts(c, listPostsBO)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, req.Page, req.PageSize)
}
