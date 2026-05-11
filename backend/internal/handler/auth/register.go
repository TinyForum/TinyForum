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
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var input request.RegisterRequest

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
