package stats

import (
	"time"
	"tiny-forum/internal/dto"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// StatsHandler 新增方法
func (h *StatsHandler) GetStatsRange(c *gin.Context) {
	var req dto.StatsRangeQuery
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
