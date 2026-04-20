package risk

// TODO: Refactory
import (
	"net/http"
	"strconv"
	"tiny-forum/internal/model"
	"tiny-forum/pkg/response"

	riskservice "tiny-forum/internal/service/risk"

	"github.com/gin-gonic/gin"
)

type RiskHandler struct {
	checkSvc *riskservice.ContentCheckService
	riskSvc  *riskservice.RiskService
}

func NewRiskHandler(checkSvc *riskservice.ContentCheckService, riskSvc *riskservice.RiskService) *RiskHandler {
	return &RiskHandler{checkSvc: checkSvc, riskSvc: riskSvc}
}

// RegisterRoutes 注册路由
func (h *RiskHandler) RegisterRoutes(admin *gin.RouterGroup) {
	g := admin.Group("/risk")
	{
		// 审核队列
		g.GET("/audit/tasks", h.ListAuditTasks)
		g.POST("/audit/tasks/:id/approve", h.ApproveTask)
		g.POST("/audit/tasks/:id/reject", h.RejectTask)

		// 操作日志
		g.GET("/audit/logs", h.ListAuditLogs)
	}
}

// ListAuditTasks 获取待审核内容列表
// GET /admin/risk/audit/tasks?limit=20&offset=0
func (h *RiskHandler) ListAuditTasks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit > 100 {
		limit = 100
	}

	tasks, total, err := h.checkSvc.GetListPendingTasks(limit, offset)
	if err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"total": total,
		"items": tasks,
	})
}

type resolveTaskInput struct {
	Note string `json:"note" binding:"max=500"`
}

// ApproveTask 审核通过（内容恢复 published）
// POST /admin/risk/audit/tasks/:id/approve
func (h *RiskHandler) ApproveTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	var input resolveTaskInput
	_ = c.ShouldBindJSON(&input)

	reviewerID := c.GetUint(model.ContextUserID)
	if err := h.checkSvc.ResolveTask(uint(taskID), true, reviewerID, input.Note); err != nil {
		response.InternalError(c, "操作失败")
		return
	}

	// 记录审计日志
	ip := c.ClientIP()
	_ = h.riskSvc.WriteAuditLog(reviewerID,
		model.AuditActionApproveContent, "audit_task", uint(taskID),
		"pending", "approved", input.Note, ip)

	response.Success(c, nil)
}

// RejectTask 审核拒绝（内容保持/改为 hidden）
// POST /admin/risk/audit/tasks/:id/reject
func (h *RiskHandler) RejectTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	var input resolveTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	reviewerID := c.GetUint(model.ContextUserID)
	if err := h.checkSvc.ResolveTask(uint(taskID), false, reviewerID, input.Note); err != nil {
		response.InternalError(c, "操作失败")
		return
	}

	ip := c.ClientIP()
	_ = h.riskSvc.WriteAuditLog(reviewerID,
		model.AuditActionRejectContent, "audit_task", uint(taskID),
		"pending", "rejected", input.Note, ip)

	response.Success(c, nil)
}

// ListAuditLogs 查询操作审计日志
// GET /admin/risk/audit/logs?target_type=post&target_id=123&limit=20
func (h *RiskHandler) ListAuditLogs(c *gin.Context) {
	targetType := c.Query("target_type")
	targetIDStr := c.Query("target_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit > 200 {
		limit = 200
	}

	var targetID uint
	if targetIDStr != "" {
		id, err := strconv.ParseUint(targetIDStr, 10, 64)
		if err != nil {
			response.BadRequest(c, "无效的 target_id")
			return
		}
		targetID = uint(id)
	}

	logs, err := h.riskSvc.GetAuditLogs(targetType, targetID, limit)
	if err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": logs})
}
