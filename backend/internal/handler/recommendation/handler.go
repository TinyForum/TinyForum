package recommendation

import (
	"strconv"

	"tiny-forum/internal/middleware"
	"tiny-forum/internal/model/request"
	recSvc "tiny-forum/internal/service/recommendation"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler 推荐系统 HTTP 处理器
type Handler struct {
	recSvc recSvc.RecommendationService
}

func NewHandler(recSvc recSvc.RecommendationService) *Handler {
	return &Handler{recSvc: recSvc}
}

// RegisterRoutes 注册推荐系统路由
func (h *Handler) RegisterRoutes(api *gin.RouterGroup, mw middleware.MiddlewareSet) {
	recGroup := api.Group("/recommendations")
	{
		// 获取个性化推荐
		recGroup.GET("/feed", mw.Auth(), h.GetFeed)

		// 记录用户行为
		recGroup.POST("/behaviors", mw.Auth(), h.RecordBehavior)
		recGroup.POST("/behaviors/batch", mw.Auth(), h.BatchRecordBehaviors)

		// 推荐反馈
		recGroup.POST("/feedback", mw.Auth(), h.SubmitFeedback)
		recGroup.POST("/feedback/batch", mw.Auth(), h.SubmitBatchFeedback)

		// 用户兴趣画像
		recGroup.GET("/profile", mw.Auth(), h.GetProfile)
		recGroup.POST("/profile/refresh", mw.Auth(), h.RefreshProfile)

		// 记录浏览（可选认证：匿名也得记录）
		recGroup.POST("/record-view", mw.OptionalAuth(), h.RecordView)
	}

	// 管理员路由
	adminRecGroup := api.Group("/admin/recommendations", mw.Auth(), mw.CasbinAuth())
	{
		adminRecGroup.GET("/overview", h.GetOverviewStats)
		adminRecGroup.GET("/behaviors", h.GetBehaviorStats)
		adminRecGroup.GET("/users", h.GetUserAnalysis)
		adminRecGroup.GET("/content", h.GetContentPerformance)
		adminRecGroup.GET("/risk-analysis", h.GetRiskAnalysis)
		adminRecGroup.GET("/user-analysis", h.GetComprehensiveUserAnalysis)

		// 用户分析 6 大模块
		adminRecGroup.GET("/ua/overview", h.GetUAOverview)
		adminRecGroup.GET("/ua/segments", h.GetUASegments)
		adminRecGroup.GET("/ua/profiles", h.GetUAProfiles)
		adminRecGroup.GET("/ua/behavior", h.GetUABehavior)
		adminRecGroup.GET("/ua/cohorts", h.GetUACohorts)
		adminRecGroup.GET("/ua/risk", h.GetUARisk)
	}
}

// GetFeed 获取个性化推荐流
// @Summary 获取个性化推荐流
// @Description 基于用户行为序列和内容特征，生成个性化推荐内容
// @Tags 推荐系统
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param strategy query string false "推荐策略" Enums(hybrid, content, collaborative, popular)
// @Success 200 {object} common.BasicResponse "推荐成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Failure 500 {object} common.BasicResponse "服务器内部错误"
// @Router /recommendations/feed [get]
func (h *Handler) GetFeed(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	strategy := c.DefaultQuery("strategy", "hybrid")

	result, err := h.recSvc.GetRecommendations(c.Request.Context(), userID, request.RecommendationQuery{
		Page:     page,
		PageSize: pageSize,
		Strategy: strategy,
	})
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// RecordBehavior 记录用户行为事件
// @Summary 记录用户行为
// @Description 记录一次用户与内容的交互行为（浏览、点赞、评论等），用于后续推荐优化
// @Tags 推荐系统
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.RecordBehaviorRequest true "行为事件"
// @Success 200 {object} common.BasicResponse "记录成功"
// @Failure 400 {object} common.BasicResponse "参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /recommendations/behaviors [post]
func (h *Handler) RecordBehavior(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.RecordBehaviorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	if err := h.recSvc.RecordBehavior(userID, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// BatchRecordBehaviors 批量记录用户行为
// @Summary 批量记录用户行为
// @Description 批量记录用户与内容的交互行为
// @Tags 推荐系统
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body []request.RecordBehaviorRequest true "行为事件列表"
// @Success 200 {object} common.BasicResponse "记录成功"
// @Failure 400 {object} common.BasicResponse "参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /recommendations/behaviors/batch [post]
func (h *Handler) BatchRecordBehaviors(c *gin.Context) {
	userID := c.GetUint("user_id")
	var events []request.RecordBehaviorRequest
	if err := c.ShouldBindJSON(&events); err != nil {
		response.HandleError(c, err)
		return
	}
	if err := h.recSvc.BatchRecordBehaviors(userID, events); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// SubmitFeedback 提交推荐反馈
// @Summary 提交推荐反馈
// @Description 用户对推荐结果的反馈（曝光、点击、停留、不感兴趣等）
// @Tags 推荐系统
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.RecommendationFeedbackRequest true "反馈信息"
// @Success 200 {object} common.BasicResponse "反馈成功"
// @Failure 400 {object} common.BasicResponse "参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /recommendations/feedback [post]
func (h *Handler) SubmitFeedback(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.RecommendationFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	if err := h.recSvc.SubmitFeedback(userID, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// SubmitBatchFeedback 批量提交推荐反馈
// @Summary 批量提交推荐反馈
// @Description 批量提交推荐反馈
// @Tags 推荐系统
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.BatchFeedbackRequest true "反馈列表"
// @Success 200 {object} common.BasicResponse "反馈成功"
// @Failure 400 {object} common.BasicResponse "参数错误"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /recommendations/feedback/batch [post]
func (h *Handler) SubmitBatchFeedback(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.BatchFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	if err := h.recSvc.SubmitBatchFeedback(userID, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// GetProfile 获取用户兴趣画像
// @Summary 获取用户兴趣画像
// @Description 获取当前用户的兴趣标签和权重分布
// @Tags 推荐系统
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse "获取成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /recommendations/profile [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	profile, err := h.recSvc.GetInterestProfile(userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, profile)
}

// RefreshProfile 刷新用户兴趣画像
// @Summary 刷新用户兴趣画像
// @Description 基于最近行为重新计算用户兴趣画像
// @Tags 推荐系统
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse "刷新成功"
// @Failure 401 {object} common.BasicResponse "未授权"
// @Router /recommendations/profile/refresh [post]
func (h *Handler) RefreshProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.recSvc.RefreshUserProfile(userID); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}

// RecordView 记录内容浏览
// @Summary 记录内容浏览
// @Description 记录用户浏览了某个内容（匿名用户也可记录，用于统计）
// @Tags 推荐系统
// @Accept json
// @Produce json
// @Param body body request.RecordBehaviorRequest true "浏览信息"
// @Success 200 {object} common.BasicResponse "记录成功"
// @Router /recommendations/record-view [post]
func (h *Handler) RecordView(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req request.RecordBehaviorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	if userID == 0 {
		userID = uint(uuid.New().ID() % 1000000) // 匿名用户用伪 ID
	}
	if req.BehaviorType == "" {
		req.BehaviorType = string("view")
	}
	if err := h.recSvc.RecordBehavior(userID, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, nil)
}
