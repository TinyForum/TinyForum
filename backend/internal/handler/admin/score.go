package admin

import (
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminGetUserScore 获取用户积分
// @Summary 获取用户积分
// @Description 获取指定用户积分，不传id则获取所有用户积分列表
// @Tags 管理员后台
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query int false "用户ID"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Failure 401 {object} common.BasicResponse
// @Failure 403 {object} common.BasicResponse
// @Failure 500 {object} common.BasicResponse
// @Router /admin/users/score [get]
func (h *AdminHandler) ListUsersScore(c *gin.Context) {
		users, err := h.service.ListUsersScore(c)
		if err != nil {
			response.InternalError(c, "查询用户积分失败")
			return
		}
		response.Success(c, users)
		return
}


// AdminGetUserScore 获取用户积分
// @Summary 获取用户积分
// @Description 获取指定用户积分，不传id则获取所有用户积分列表
// @Tags 管理员后台
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query int false "用户ID"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Failure 401 {object} common.BasicResponse
// @Failure 403 {object} common.BasicResponse
// @Failure 500 {object} common.BasicResponse
// @Router /admin/users/score [get]
func (h *AdminHandler) GetUserScore(c *gin.Context) {
var 	req request.GetUserScoreRequest
	err := c.BindQuery(&req)
	if err != nil {
	    response.HandleError(c, err)
		return
	}

		scoreVO, err := h.service.GetUserScore(c,req.UserID)
		if err != nil {
			response.HandleError(c, err)
			return
		}
		response.Success(c, scoreVO)
	
}