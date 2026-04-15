package handler

import (
	"fmt"
	"time"
	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"
	"tiny-forum/pkg/utils"

	"github.com/gin-gonic/gin"
)

// StatsHandler 统计 HTTP 处理器
type StatsHandler struct {
	statsSvc *service.StatsService
}

func NewStatsHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{statsSvc: svc}
}

// GetStatsDay 获取日统计数据
// @Summary 获取日统计数据
// @Description 获取指定日期的统计数据（用户、帖子、评论等）
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param date query string false "日期 (格式: 2006-01-02)" default(今天)
// @Param type query string false "统计类型" Enums(users, posts, comments, all) default(all)
// @Success 200 {object} response.Response{data=model.StatsTodayInfo}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/statistics/day [get]
func (h *StatsHandler) GetStatsDay(c *gin.Context) {
	var req struct {
		Date string `form:"date" binding:"omitempty,datetime=2006-01-02"`
		Type string `form:"type" binding:"omitempty,oneof=users posts comments all"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	// if req.Date == "" {
	// 	req.Date = time.Now().Format("2006-01-02")
	// }
	date, err := utils.ParseTimeExpression(req.Date, time.Now(), time.Local, false)
	if req.Type == "" {
		req.Type = "all"
	}

	stats, err := h.statsSvc.GetStatsByDate(c.Request.Context(), date, req.Type)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "获取统计数据失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "data": stats})
}

// GetStatsTotal 获取总计统计数据
// @Summary 获取总计统计数据
// @Description 获取指定时间范围内的汇总统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param start_date query string false "开始日期 (格式: 2006-01-02，默认30天前)"
// @Param end_date   query string false "结束日期 (格式: 2006-01-02，默认今天)"
// @Param type       query string false "统计类型" Enums(users, posts, comments, all) default(all)
// @Success 200 {object} response.Response{data=model.StatsInfoResp}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/statistics/total [get]
func (h *StatsHandler) GetStatsTotal(c *gin.Context) {
	var req struct {
		StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
		EndDate   string `form:"end_date"   binding:"omitempty,datetime=2006-01-02"`
		Type      string `form:"type"       binding:"omitempty,oneof=users posts comments all"`
	}

	// ✅ 正确顺序：先绑定参数
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "invalid parameters: "+err.Error())
		return
	}

	// 设置默认值
	if req.Type == "" {
		req.Type = "all"
	}

	// ✅ 然后解析时间范围（此时 req 已经有值了）
	date, err := utils.ParseTimeRange(req.StartDate, req.EndDate)
	if err != nil {
		response.BadRequest(c, "invalid date range: "+err.Error())
		return
	}

	fmt.Printf("Parsed time range: %v to %v\n", date.Start, date.End)
	fmt.Printf("User requested: start_date=%s, end_date=%s, type=%s\n",
		req.StartDate, req.EndDate, req.Type)

	// 调用 Service
	totals, err := h.statsSvc.GetTotalStats(c.Request.Context(), date.Start, date.End, req.Type)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "获取总计数据失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 0, "data": totals})
}

// GetStatsTrend 获取趋势统计数据
// @Summary 获取趋势统计数据
// @Description 获取指定时间范围内的趋势数据（按天、周、月统计）
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param start_date query string false "开始日期 (格式: 2006-01-02，默认7天前)"
// @Param end_date   query string false "结束日期 (格式: 2006-01-02，默认今天)"
// @Param type       query string true  "统计类型" Enums(users, posts, comments)
// @Param interval   query string false "统计粒度" Enums(day, week, month) default(day)
// @Success 200 {object} response.Response{data=model.StatsTrendResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/statistics/trend [get]
func (h *StatsHandler) GetStatsTrend(c *gin.Context) {
	var req struct {
		StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
		EndDate   string `form:"end_date"   binding:"omitempty,datetime=2006-01-02"`
		Type      string `form:"type"       binding:"required,oneof=users posts comments"`
		Interval  string `form:"interval"   binding:"omitempty,oneof=day week month"`
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	date, err := utils.ParseTimeRange(req.StartDate, req.EndDate)
	if req.Interval == "" {
		req.Interval = "day"
	}

	// Service 接收字符串日期，内部自行解析
	trend, err := h.statsSvc.GetTrendStats(c.Request.Context(), date.Start, date.End, req.Type, req.Interval)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "message": "获取趋势数据失败: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 0,
		"data": gin.H{
			"start_date": date.Start.Format("2026-01-02"),
			"end_date":   date.End.Format("2026-01-02"),
			"interval":   req.Interval,
			"type":       req.Type,
			"trend":      trend,
		},
	})
}

// ── 辅助函数 ──────────────────────────────────────────────────────────────────

// parseDateRangeWithDefault 将可选的日期字符串解析为 "YYYY-MM-DD" 格式。
// 若两者均为空则以当前时间为终点，向前推 defaultDays 天为起点。
func parseDateRangeWithDefault(startDate, endDate string, defaultDays int) (string, string) {
	now := time.Now()

	var start, end time.Time

	if endDate != "" {
		end, _ = time.ParseInLocation("2006-01-02", endDate, time.Local)
	} else {
		end = now
	}

	if startDate != "" {
		start, _ = time.ParseInLocation("2006-01-02", startDate, time.Local)
	} else {
		start = end.AddDate(0, 0, -defaultDays)
	}

	// 防止 start > end
	if start.After(end) {
		start = end.AddDate(0, 0, -defaultDays)
	}

	return start.Format("2006-01-02"), end.Format("2006-01-02")
}
