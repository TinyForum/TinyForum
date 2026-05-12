package admin

import (
	"strconv"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *AdminHandler) ListReviewRequire(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	listPostsBO := &common.PageQuery[bo.ListPosts]{
		Page:     page,
		PageSize: pageSize,
		Data: bo.ListPosts{
			PostStatus: do.PostStatusPending,
			Keyword:    keyword,
		},
	}
	posts, total, err := h.service.ListReviewRequire(c, listPostsBO)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}
