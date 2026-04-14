package handler

import (
	"time"
	"tiny-forum/internal/service"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsSvc *service.StatsService
}

func NewStatsHandler(service *service.StatsService) *StatsHandler {
	return &StatsHandler{
		statsSvc: service,
	}
}

// GetStatsDay 获取日统计数据
// @Summary 获取日统计数据
// @Description 获取指定日期的统计数据（用户、帖子、评论等）
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param date query string false "日期 (格式: 2006-01-02)" default(今天)
// @Param type query string false "统计类型" Enums(users, posts, comments, likes, all) default(all)
// @Success 200 {object} response.Response{data=model.StatsTodayInfo} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/statistics/day [get]
func (h *StatsHandler) GetStatsDay(c *gin.Context) {
	var req struct {
		Date string `form:"date" binding:"omitempty,datetime=2006-01-02"`
		Type string `form:"type" binding:"omitempty,oneof=users posts comments likes all"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 默认今天
	if req.Date == "" {
		req.Date = time.Now().Format("2006-01-02")
	}

	// 获取指定日期的统计数据
	stats, err := h.statsSvc.GetStatsByDate(c.Request.Context(), req.Date, req.Type)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "获取统计数据失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": stats,
	})
}

// GetStatsTotal 获取总计统计数据
// @Summary 获取总计统计数据
// @Description 获取指定时间范围内的总计统计数据（支持用户、帖子、评论等类型）
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param start_date query string false "开始日期 (格式: 2006-01-02)" default(30天前)
// @Param end_date query string false "结束日期 (格式: 2006-01-02)" default(今天)
// @Param type query string false "统计类型" Enums(users, posts, comments, likes, all) default(all)
// @Success 200 {object} response.Response{data=model.StatsTotalResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/statistics/total [get]
func (h *StatsHandler) GetStatsTotal(c *gin.Context) {
	var req struct {
		StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
		EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
		Type      string `form:"type" binding:"omitempty,oneof=users posts comments likes all"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 处理日期范围
	start, end := parseDateRangeWithDefault(req.StartDate, req.EndDate, 30) // 默认30天

	// 获取总计数据
	totals, err := h.statsSvc.GetTotalStats(c.Request.Context(), start, end, req.Type)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "获取总计数据失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"start_date": start,
			"end_date":   end,
			"totals":     totals,
		},
	})
}

func parseDateRangeWithDefault(startDate, endDate string, defaultDays int) (string, string) {
	now := time.Now()

	if startDate == "" && endDate == "" {
		endDate = now.Format("2006-01-02")
		startDate = now.AddDate(0, 0, -defaultDays).Format("2006-01-02")
		return startDate, endDate
	}

	if startDate == "" {
		startDate = now.AddDate(0, 0, -defaultDays).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = now.Format("2006-01-02")
	}

	return startDate, endDate
}

// GetStatsTrend 获取趋势统计数据
// @Summary 获取趋势统计数据
// @Description 获取指定时间范围内的趋势数据（按天、周、月统计）
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param start_date query string false "开始日期 (格式: 2006-01-02)" default(7天前)
// @Param end_date query string false "结束日期 (格式: 2006-01-02)" default(今天)
// @Param type query string true "统计类型" Enums(users, posts, comments, likes)
// @Param interval query string false "统计粒度" Enums(day, week, month) default(day)
// @Success 200 {object} response.Response{data=model.StatsTrendResponse} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/statistics/trend [get]
func (h *StatsHandler) GetStatsTrend(c *gin.Context) {
	var req struct {
		StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
		EndDate   string `form:"end_date" binding:"omitempty,datetime=2006-01-02"`
		Type      string `form:"type" binding:"required,oneof=users posts comments likes"`
		Interval  string `form:"interval" binding:"omitempty,oneof=day week month"` // 统计粒度
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 默认区间：最近7天
	start, end := parseDateRangeWithDefault(req.StartDate, req.EndDate, 7)

	// 默认粒度：day
	if req.Interval == "" {
		req.Interval = "day"
	}

	// 获取趋势数据
	trend, err := h.statsSvc.GetTrendStats(c.Request.Context(), start, end, req.Type, req.Interval)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "获取趋势数据失败"})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"start_date": start,
			"end_date":   end,
			"interval":   req.Interval,
			"type":       req.Type,
			"trend":      trend,
		},
	})
}
