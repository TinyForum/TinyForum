package auth

import (
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
// @Summary 用户注册
// @Tags 验证管理
// @Accept json
// @Produce json
// @Param body body request.RegisterRequest true "注册信息"
// @Success 200 {object} vo.UserPrivateVO
// @Failure 400 {object} common.BasicResponse
// @Router /auth/register [post]
// Register 是一个处理用户注册请求的方法
// 它接收一个 gin.Context 对象作为参数，用于处理 HTTP 请求和响应
func (h *AuthHandler) Register(c *gin.Context) {
	// 从请求上下文中获取 context
	ctx := c.Request.Context()
	// 声明一个 RegisterRequest 类型的变量 input，用于存储注册请求的数据
	var input request.RegisterRequest

	// 将请求体中的 JSON 数据绑定到 input 结构体上
	// 如果绑定失败，则处理错误并返回
	if err := c.ShouldBindJSON(&input); err != nil {
		response.HandleError(c, err)
		return
	}

	result, err := h.authSvc.Register(ctx, input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, result)
}
