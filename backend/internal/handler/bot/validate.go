package bot

import (
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ValidateFlow 校验零代码 Flow 配置
// @Summary 校验零代码流程
// @Tags 机器人管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body ValidateFlowRequest true "Flow 配置"
// @Success 200 {object} common.BasicResponse{data=object{valid=bool,errors=array}}
// @Router /bots/nocode/validate [post]
func (h *Handler) ValidateFlow(c *gin.Context) {
	var req request.ValidateFlowRequest
	// 打印请求
	if err := c.ShouldBindJSON(&req); err != nil {

		response.HandleError(c, err)
		return
	}
	flow := req.ToFlow()
	errs := h.svc.ValidateFlow(&flow)
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		msgs = append(msgs, e.Error())
	}
	response.Success(c, gin.H{
		"valid":  len(errs) == 0,
		"errors": msgs,
	})
}
