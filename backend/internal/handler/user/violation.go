package user

import (
	"strconv"

	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
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
// @Router /users/me/violations [get]
func (h *UserHandler) ListUserViolations(c *gin.Context) {
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
}

// GetUserViolation 获取当前用户的违规记录详情
// @Summary 获取违规详情
// @Description 获取当前登录用户的单条违规记录详情
// @Tags 用户中心
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "违规记录ID"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Failure 404 {object} common.BasicResponse
// @Router /users/me/violations/{id} [get]
func (h *UserHandler) GetUserViolation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}
	userID := c.GetUint("user_id")
	result, err := h.userSvc.GetUserViolation(c.Request.Context(), uint(id), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// AppealUserViolation 提交违规申诉
// @Summary 提交违规申诉
// @Description 当前登录用户对指定违规记录提交申诉
// @Tags 用户中心
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "违规记录ID"
// @Param body body request.AppealViolationRequest true "申诉请求"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Failure 404 {object} common.BasicResponse
// @Router /users/me/violations/{id}/appeal [post]
func (h *UserHandler) AppealUserViolation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}
	var req request.AppealViolationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	userID := c.GetUint("user_id")
	result, err := h.userSvc.AppealUserViolation(c.Request.Context(), uint(id), userID, req.Reason)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}
