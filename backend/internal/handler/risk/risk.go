package risk

import (
	"net/http"
	"strconv"
	"tiny-forum/internal/model/do"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListAuditTasks
// @Summary      获取待审核内容列表
// @Description  分页查询所有状态为 pending 的审核任务（内容安全检测不通过待人工审核）
// @Tags         风险管理
// @Accept       json
// @Produce      json
// @Param        limit   query     int     false  "每页数量，默认20，最大100"  default(20)
// @Param        offset  query     int     false  "偏移量，默认0"            default(0)
// @Success      200     {object}  vo.BasicResponse "成功"
// @Failure      500     {object}  vo.BasicResponse "服务器内部错误"
// @Router       /admin/risk/audit/tasks [get]
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

// ApproveTask
// @Summary      审核通过
// @Description  将指定审核任务标记为通过，恢复内容为 published 状态，并记录操作日志
// @Tags         风险管理
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "审核任务ID"
// @Param        input body      resolveTaskInput  false "审核备注"
// @Success      200   {object}  vo.BasicResponse  "操作成功"
// @Failure      400   {object}  vo.BasicResponse "请求参数错误"
// @Failure      500   {object}  vo.BasicResponse "服务器内部错误"
// @Router       /admin/risk/audit/tasks/{id}/approve [post]
func (h *RiskHandler) ApproveTask(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}

	var input resolveTaskInput
	_ = c.ShouldBindJSON(&input)

	reviewerID := c.GetUint(do.ContextUserID)
	if err := h.checkSvc.ResolveTask(uint(taskID), true, reviewerID, input.Note); err != nil {
		response.InternalError(c, "操作失败")
		return
	}

	// 记录审计日志
	ip := c.ClientIP()
	_ = h.riskSvc.WriteAuditLog(reviewerID,
		do.AuditActionApproveContent, "audit_task", uint(taskID),
		"pending", "approved", input.Note, ip)

	response.Success(c, nil)
}

// RejectTask
// @Summary      审核拒绝
// @Description  将指定审核任务标记为拒绝，内容状态改为 hidden（屏蔽），并记录操作日志
// @Tags         风险管理
// @Accept       json
// @Produce      json
// @Param        id    path      int               true  "审核任务ID"
// @Param        input body      resolveTaskInput  true  "审核备注"
// @Success      200   {object}  vo.BasicResponse  "操作成功"
// @Failure      400   {object}  vo.BasicResponse "请求参数错误"
// @Failure      500   {object}  vo.BasicResponse "服务器内部错误"
// @Router       /admin/risk/audit/tasks/{id}/reject [post]
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

	reviewerID := c.GetUint(do.ContextUserID)
	if err := h.checkSvc.ResolveTask(uint(taskID), false, reviewerID, input.Note); err != nil {
		response.InternalError(c, "操作失败")
		return
	}

	ip := c.ClientIP()
	_ = h.riskSvc.WriteAuditLog(reviewerID,
		do.AuditActionRejectContent, "audit_task", uint(taskID),
		"pending", "rejected", input.Note, ip)

	response.Success(c, nil)
}

// ListAuditLogs
// @Summary      查询审核操作日志
// @Description  根据目标类型和ID查询管理员对审核任务的操作记录（批准/拒绝）
// @Tags         风险管理
// @Accept       json
// @Produce      json
// @Param        target_type  query     string  false  "目标类型，如 post、comment、user 等"
// @Param        target_id    query     int     false  "目标ID"
// @Param        limit        query     int     false  "返回条数，默认50，最大200"  default(50)
// @Success      200     {object}  vo.BasicResponse "成功"
// @Failure      400          {object}  vo.BasicResponse "参数错误（target_id无效）"
// @Failure      500          {object}  vo.BasicResponse "服务器内部错误"
// @Router       /admin/risk/audit/logs [get]
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
