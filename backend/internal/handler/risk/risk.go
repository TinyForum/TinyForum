package risk

import (
	"strconv"
	"tiny-forum/internal/model/vo"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListAuditLogs
// @Summary      查询审核操作日志
// @Description  根据目标类型和ID查询管理员对审核任务的操作记录（批准/拒绝）
// @Tags         风险管理
// @Accept       json
// @Produce      json
// @Param        target_type  query     string  false  "目标类型，如 post、comment、user 等"
// @Param        target_id    query     int     false  "目标ID"
// @Param        limit        query     int     false  "返回条数，默认50，最大200"  default(50)
// @Success      200     {object}  common.BasicResponse "成功"
// @Failure      400          {object}  common.BasicResponse "参数错误（target_id无效）"
// @Failure      500          {object}  common.BasicResponse "服务器内部错误"
// @Router       /risk/audit/logs [get]
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
			response.HandleError(c, err)
			return
		}
		targetID = uint(id)
	}

	logs, err := h.riskSvc.GetAuditLogs(targetType, targetID, limit)
	if err != nil {
		response.InternalError(c, "查询失败")
		return
	}

	responseData := vo.AuditLogsVO{
		Log: logs,
	}
	response.Success(c, responseData)
}
