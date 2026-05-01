package user

import (
	"fmt"
	"tiny-forum/internal/model/do"

	"github.com/gin-gonic/gin"
)

// getViewerID 从 context 获取当前登录用户 ID，未登录返回 0
func getViewerID(c *gin.Context) uint {
	if v, exists := c.Get("user_id"); exists {
		if id, ok := v.(uint); ok {
			return id
		}
	}
	return 0
}

// handleRoleError 统一处理角色变更错误

// sendTempPasswordNotification 发送临时密码通知（内部辅助）
func (h *UserHandler) sendTempPasswordNotification(targetID, operatorID uint, tempPassword string) {
	message := fmt.Sprintf(
		"管理员已重置您的密码。临时密码为：%s，有效期 30 分钟，请尽快登录并修改密码，以防被盗。",
		tempPassword,
	)
	h.notifSvc.Create(targetID, &operatorID, do.NotifySystem, message, nil, "")
}

// ── 共享请求/响应结构体 ────────────────────────────────────────────────

// LeaderboardRequest 排行榜请求参数

// LeaderboardResponse 排行榜响应

// AdminSetScoreRequest 管理员设置积分请求
type AdminSetScoreRequest struct {
	Operation string `json:"operation" binding:"required,oneof=set add subtract"`
	Score     int    `json:"score" binding:"required,gte=0,lte=999999"`
	Reason    string `json:"reason" binding:"required,max=200"`
}

// AdminSetScoreResponse 管理员设置积分响应
type AdminSetScoreResponse struct {
	UserID     uint64 `json:"user_id"`
	OldScore   int    `json:"old_score"`
	NewScore   int    `json:"new_score"`
	Change     int    `json:"change"`
	Operation  string `json:"operation"`
	OperatorID uint   `json:"operator_id"`
	Reason     string `json:"reason"`
	Timestamp  int64  `json:"timestamp"`
}

// AdminResetUserPasswordResponse 重置密码响应
type AdminResetUserPasswordResponse struct {
	Message    string `json:"message"`
	UserID     uint   `json:"user_id"`
	OperatorID uint   `json:"operator_id"`
}
