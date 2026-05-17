package bot

import (
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ─── 零代码 ───────────────────────────────────────────────────────────────

// GetNocodeMetadata 获取零代码编辑器所需的内置节点元数据
// @Summary 获取零代码节点元数据
// @Description 返回所有内置 Trigger/Condition/Action 的类型和参数定义，供前端拖拽编辑器使用
// @Tags 机器人管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse{data=nocode.NocodeMetadata}
// @Router /bots/nocode/metadata [get]
func (h *Handler) GetNocodeMetadata(c *gin.Context) {
	response.Success(c, h.svc.GetNocodeMetadata())
}
