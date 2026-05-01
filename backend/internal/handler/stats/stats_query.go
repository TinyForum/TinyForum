package stats

import (
	"fmt"
	"log"
	"time"
	"tiny-forum/internal/model/dto"
	"tiny-forum/pkg/response"

	// "tiny-forum/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetStatsDay 获取日统计数据
// @Summary 获取日统计数据
// @Description 获取指定日期的统计数据（用户、帖子、评论等）
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param date query string false "日期 (格式: 2006-01-02)" default(今天)
// @Param type query string false "统计类型" Enums(users, posts, comments, all) default(all)
// @Success 200 {object} response.Response{data=po.StatsTodayInfo}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /statistics/day [get]
func (h *StatsHandler) GetStatsDay(c *gin.Context) {
	var req dto.StatsDayQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Type == "" {
		req.Type = "all"
	}

	dateStr := req.Date
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := h.timeHelpers.SingleParser.Parse(dateStr, time.Now(), time.Local, false)
	if err != nil {
		response.BadRequest(c, "无效的日期格式: "+err.Error())
		return
	}

	stats, err := h.statsSvc.GetStatsByDate(c.Request.Context(), date, req.Type)
	if err != nil {
		// 记录详细错误日志，便于排查
		log.Printf("获取统计数据失败: %v", err) // 或者使用 slog
		response.InternalError(c, "获取统计数据失败，请稍后重试")
		return
	}

	response.Success(c, stats)
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
// @Success 200 {object} response.Response{data=po.StatsInfoResp}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /statistics/total [get]
func (h *StatsHandler) GetStatsTotal(c *gin.Context) {
	var req struct {
		StartDate string `form:"start_date" binding:"omitempty,datetime=2006-01-02"`
		EndDate   string `form:"end_date"   binding:"omitempty,datetime=2006-01-02"`
		Type      string `form:"type"       binding:"omitempty,oneof=users posts comments all"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "invalid parameters: "+err.Error())
		return
	}

	if req.Type == "" {
		req.Type = "all"
	}

	date, err := h.timeHelpers.RangeParser.Parse(req.StartDate, req.EndDate)
	if err != nil {
		response.BadRequest(c, "invalid date range: "+err.Error())
		return
	}

	fmt.Printf("Parsed time range: %v to %v\n", date.Start, date.End)
	fmt.Printf("User requested: start_date=%s, end_date=%s, type=%s\n",
		req.StartDate, req.EndDate, req.Type)

	totals, err := h.statsSvc.GetTotalStats(c.Request.Context(), date.Start, date.End, req.Type)
	if err != nil {
		response.InternalError(c, "获取总计数据失败: "+err.Error())
		return
	}

	response.Success(c, totals)
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
// @Success 200 {object} response.Response{data=po.StatsTrendResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /statistics/trend [get]
func (h *StatsHandler) GetStatsTrend(c *gin.Context) {
	var req dto.AdminStatsTrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	date, err := h.timeHelpers.RangeParser.Parse(req.StartDate, req.EndDate)
	if err != nil {
		response.BadRequest(c, "invalid date range: "+err.Error())
		return
	}
	if req.Interval == "" {
		req.Interval = "day"
	}

	trend, err := h.statsSvc.GetTrendStats(c.Request.Context(), date.Start, date.End, req.Type, req.Interval)
	if err != nil {
		response.InternalError(c, "获取趋势数据失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"start_date": date.Start.Format("2006-01-02"),
		"end_date":   date.End.Format("2006-01-02"),
		"interval":   req.Interval,
		"type":       req.Type,
		"trend":      trend,
	})
}
