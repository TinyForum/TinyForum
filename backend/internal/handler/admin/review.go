package admin

import (
	"strconv"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListReviewRequire 获取待审核帖子列表
// @Summary 获取待审核帖子列表
// @Tags 帖子管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "关键词"
// @Success 200 {object} common.BasicResponse "获取成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "无权限（需要管理员权限）"
// @Failure 500 {object} common.BasicResponse "服务器内部错误"
// @Router /admin/posts/pending [get]
func (h *AdminHandler) ListReviewRequire(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	listPostsBO := &common.PageQuery[bo.ListPosts]{
		Page:     page,
		PageSize: pageSize,
		Data: bo.ListPosts{
			PostStatus: do.CreationStatusPending,
			Keyword:    keyword,
		},
	}
	posts, total, err := h.service.ListReviewRequire(c, listPostsBO)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}
