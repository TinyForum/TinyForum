package timeline

import (
	"strconv"

	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetHomeTimeline 获取首页时间线
// @Summary 获取首页时间线
// @Description 获取当前用户首页的时间线（推荐内容）
// @Tags 时间线管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]do.TimelineEvent}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /timeline/home [get]
func (h *TimelineHandler) GetHomeTimeline(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	events, total, err := h.timelineSvc.GetHomeTimeline(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, events, total, page, pageSize)
}

// GetFollowingTimeline 获取关注时间线
// @Summary 获取关注时间线
// @Description 获取当前用户关注的人的内容时间线
// @Tags 时间线管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]do.TimelineEvent}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /timeline/following [get]
func (h *TimelineHandler) GetFollowingTimeline(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	events, total, err := h.timelineSvc.GetFollowingTimeline(userID, page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, events, total, page, pageSize)
}
