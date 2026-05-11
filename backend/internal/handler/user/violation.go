package user

import (
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// @Summary 获取用户违规记录
// @Description 获取用户违规记录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /users/me/violation [get]
func (h *UserHandler) ListUserViolation(c *gin.Context) {
	// TODO implement me
	userID := c.GetUint("user_id")

	var req request.ListUserViolationRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}

	result, err := h.userSvc.ListUserViolation(c.Request.Context(), req, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, result)
	panic("implement me")
}
