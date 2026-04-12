package handler

import (
	"strconv"

	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type TimelineHandler struct {
	timelineSvc *service.TimelineService
}

func NewTimelineHandler(timelineSvc *service.TimelineService) *TimelineHandler {
	return &TimelineHandler{timelineSvc: timelineSvc}
}

// GetHomeTimeline 获取首页时间线
// @Summary 获取首页时间线
// @Description 获取当前用户首页的时间线（推荐内容）
// @Tags 时间线管理
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.TimelineEvent}} "获取成功"
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
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.TimelineEvent}} "获取成功"
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

// Subscribe 关注用户
// @Summary 关注用户
// @Description 关注指定用户，接收其动态更新
// @Tags 时间线管理
// @Produce json
// @Security ApiKeyAuth
// @Param user_id path int true "要关注的用户ID"
// @Success 200 {object} response.Response{data=object} "关注成功"
// @Failure 400 {object} response.Response "无效的用户ID或不能关注自己"
// @Failure 401 {object} response.Response "未授权"
// @Failure 409 {object} response.Response "已关注该用户"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /timeline/subscribe/{user_id} [post]
func (h *TimelineHandler) Subscribe(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	subscriberID := c.GetUint("user_id")
	if subscriberID == uint(userID) {
		response.BadRequest(c, "不能关注自己")
		return
	}

	if err := h.timelineSvc.Subscribe(subscriberID, uint(userID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "关注成功"})
}

// Unsubscribe 取消关注用户
// @Summary 取消关注用户
// @Description 取消关注指定用户
// @Tags 时间线管理
// @Produce json
// @Security ApiKeyAuth
// @Param user_id path int true "要取消关注的用户ID"
// @Success 200 {object} response.Response{data=object} "取消关注成功"
// @Failure 400 {object} response.Response "无效的用户ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "未关注该用户"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /timeline/subscribe/{user_id} [delete]
func (h *TimelineHandler) Unsubscribe(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	subscriberID := c.GetUint("user_id")

	if err := h.timelineSvc.Unsubscribe(subscriberID, uint(userID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "取消关注成功"})
}

// GetSubscriptions 获取关注列表
// @Summary 获取关注列表
// @Description 获取当前用户关注的所有用户列表
// @Tags 时间线管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]model.User} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /timeline/subscriptions [get]
func (h *TimelineHandler) GetSubscriptions(c *gin.Context) {
	userID := c.GetUint("user_id")

	subscriptions, err := h.timelineSvc.GetSubscriptions(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, subscriptions)
}

// IsSubscribed 检查是否已关注
// @Summary 检查是否已关注
// @Description 检查当前用户是否已关注指定用户
// @Tags 时间线管理
// @Produce json
// @Security ApiKeyAuth
// @Param user_id path int true "要检查的用户ID"
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "无效的用户ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /timeline/subscribe/{user_id}/status [get]
func (h *TimelineHandler) IsSubscribed(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	subscriberID := c.GetUint("user_id")

	isSubscribed, err := h.timelineSvc.IsSubscribed(subscriberID, uint(userID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"is_subscribed": isSubscribed})
}
