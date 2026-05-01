package user

import (
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetStatistics 获取全局统计数据（总帖子数、总评论数、总收藏数、当前用户未读通知数）
// @Summary 获取全局统计数据
// @Tags 统计
// @Security BearerAuth
// @Success 200 {object} vo.Statistics "返回统计数据"
// @Router /statistics [get]
func (h *UserHandler) GetStatisticsCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := h.userSvc.GetGlobalStatsCount(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, stats)
}
