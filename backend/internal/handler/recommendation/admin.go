package recommendation

import (
	"strconv"

	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetOverviewStats 获取推荐系统概览统计
// @Summary 管理员：推荐系统概览
// @Description 获取推荐系统的核心统计数据（行为总数、反馈率、点击率等）
// @Tags 管理员-推荐系统
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse "成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "权限不足"
// @Router /admin/recommendations/overview [get]
func (h *Handler) GetOverviewStats(c *gin.Context) {
	stats, err := h.recSvc.GetOverviewStats(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, stats)
}

// GetBehaviorStats 获取行为统计数据
// @Summary 管理员：行为统计分析
// @Description 获取用户行为分布、每日趋势和热门目标内容
// @Tags 管理员-推荐系统
// @Produce json
// @Security ApiKeyAuth
// @Param days query int false "统计天数" default(7)
// @Success 200 {object} common.BasicResponse "成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "权限不足"
// @Router /admin/recommendations/behaviors [get]
func (h *Handler) GetBehaviorStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	stats, err := h.recSvc.GetBehaviorStats(c.Request.Context(), days)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, stats)
}

// GetUserAnalysis 获取用户分析
// @Summary 管理员：用户行为分析
// @Description 获取活跃用户排行、行为画像和热门兴趣标签
// @Tags 管理员-推荐系统
// @Produce json
// @Security ApiKeyAuth
// @Param days query int false "统计天数" default(7)
// @Success 200 {object} common.BasicResponse "成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "权限不足"
// @Router /admin/recommendations/users [get]
func (h *Handler) GetUserAnalysis(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "7"))
	analysis, err := h.recSvc.GetUserAnalysis(c.Request.Context(), days)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, analysis)
}

// GetContentPerformance 获取内容表现
// @Summary 管理员：内容表现分析
// @Description 获取热度/质量最高内容、板块分布和质量分分布
// @Tags 管理员-推荐系统
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse "成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "权限不足"
// @Router /admin/recommendations/content [get]
func (h *Handler) GetContentPerformance(c *gin.Context) {
	perf, err := h.recSvc.GetContentPerformance(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, perf)
}

// GetRiskAnalysis 获取风控关联分析
// @Summary 管理员：风控关联分析
// @Description 结合推荐行为数据与风险/违规数据，展示风控用户的关联分析
// @Tags 管理员-推荐系统
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse "成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "权限不足"
// @Router /admin/recommendations/risk-analysis [get]
func (h *Handler) GetRiskAnalysis(c *gin.Context) {
	analysis, err := h.recSvc.GetRiskAnalysis(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, analysis)
}

// GetComprehensiveUserAnalysis 获取用户综合分析
// @Summary 管理员：用户综合分析
// @Description 多维度用户分析：画像、标签分布、行为模式、相似用户群组、风险画像
// @Tags 管理员-推荐系统
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse "成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 403 {object} common.BasicResponse "权限不足"
// @Router /admin/recommendations/user-analysis [get]
func (h *Handler) GetComprehensiveUserAnalysis(c *gin.Context) {
	analysis, err := h.recSvc.GetComprehensiveUserAnalysis(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, analysis)
}
