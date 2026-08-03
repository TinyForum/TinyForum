package recommendation

import (
	"strconv"

	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ── 1. 用户概览 ──

// GetUAOverview 管理员：用户分析概览
// @Summary 管理员：用户分析概览
// @Tags 管理员-用户分析
// @Router /admin/recommendations/ua/overview [get]
func (h *Handler) GetUAOverview(c *gin.Context) {
	data, err := h.recSvc.GetUAOverview(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, data)
}

// ── 2. 用户分群 ──

// GetUASegments 管理员：用户分群
// @Tags 管理员-用户分析
// @Router /admin/recommendations/ua/segments [get]
func (h *Handler) GetUASegments(c *gin.Context) {
	data, err := h.recSvc.GetUASegments(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, data)
}

// ── 2b. 用户画像列表 ──

// GetUAProfiles 管理员：用户画像列表
// @Tags 管理员-用户分析
// @Router /admin/recommendations/ua/profiles [get]
func (h *Handler) GetUAProfiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	tier := c.Query("tier")
	sortBy := c.Query("sort_by")
	data, err := h.recSvc.GetUAProfiles(c.Request.Context(), page, pageSize, keyword, tier, sortBy)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, data)
}

// ── 3. 行为分析 ──

// GetUABehavior 管理员：行为分析
// @Tags 管理员-用户分析
// @Router /admin/recommendations/ua/behavior [get]
func (h *Handler) GetUABehavior(c *gin.Context) {
	data, err := h.recSvc.GetUABehavior(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, data)
}

// ── 4. 同期群 ──

// GetUACohorts 管理员：同期群分析
// @Tags 管理员-用户分析
// @Router /admin/recommendations/ua/cohorts [get]
func (h *Handler) GetUACohorts(c *gin.Context) {
	data, err := h.recSvc.GetUACohorts(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, data)
}

// ── 5. 风险评估 ──

// GetUARisk 管理员：用户风险评估
// @Tags 管理员-用户分析
// @Router /admin/recommendations/ua/risk [get]
func (h *Handler) GetUARisk(c *gin.Context) {
	data, err := h.recSvc.GetUARisk(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, data)
}
