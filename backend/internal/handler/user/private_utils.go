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

// AdminSetScoreRequest 管理员设置积分请求
