package stats

import (
	"time"
	"tiny-forum/internal/model/dto"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary 获取指定日期范围的统计数据
// @Description 根据日期范围和统计类型获取统计数据
// @Tags 统计管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param type query string false "统计类型" Enums(all, user, post, comment, like, board) default(all)
// @Param start_date query string false "开始日期，格式：2006-01-02"
// @Param end_date query string false "结束日期，格式：2006-01-02"
// @Success 200 {object} response.Response{data=object} "获取成功"
// @Failure 400 {object} response.Response "参数错误或日期范围超过90天"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /statistics/range [get]
func (h *StatsHandler) GetStatsRange(c *gin.Context) {
	var req dto.StatsRangeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 默认类型
	if req.Type == "" {
		req.Type = "all"
	}

	// 解析日期范围，默认最近30天
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -29) // 最近30天（包含今天）

	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			response.BadRequest(c, "结束日期格式错误: "+err.Error())
			return
		}
		endDate = parsed
	}
	if req.StartDate != "" {
		parsed, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			response.BadRequest(c, "开始日期格式错误: "+err.Error())
			return
		}
		startDate = parsed
	}

	// 校验范围不超过90天（可选）
	if endDate.Sub(startDate).Hours()/24 > 90 {
		response.BadRequest(c, "日期范围不能超过90天")
		return
	}

	stats, err := h.statsSvc.GetStatsByDateRange(c.Request.Context(), startDate, endDate, req.Type)
	if err != nil {
		response.InternalError(c, "获取统计数据失败: "+err.Error())
		return
	}

	response.Success(c, stats)
}
