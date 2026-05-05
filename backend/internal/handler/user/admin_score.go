package user

import (
	"strconv"
	"time"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminSetScore 管理员设置用户积分
// @Summary 管理员设置用户积分
// @Description 管理员可以通过此接口对指定用户的积分进行设置、增加或扣除操作。支持三种操作模式：<br>
// @Description - **set**：将用户积分设置为指定值<br>
// @Description - **add**：在现有积分基础上增加指定分数<br>
// @Description - **subtract**：从现有积分中扣除指定分数<br>
// @Description 积分范围限制为 0 ~ 999999，且操作后积分不能为负数或超出上限。
// @Tags 管理员后台
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID" example(10086)
// @Param body body AdminSetScoreRequest true "积分操作信息"
// @Success 200 {object} vo.BasicResponse "操作成功"
// @Failure 400 {object} vo.BasicResponse"请求参数错误（如积分范围非法、操作类型错误等）"
// @Failure 401 {object} vo.BasicResponse"未授权（缺少或无效的认证令牌）"
// @Failure 403 {object} vo.BasicResponse"禁止访问（当前管理员无权限操作该用户）"
// @Failure 500 {object} vo.BasicResponse"服务器内部错误（如数据库操作失败）"
// @Router /admin/users/{id}/score [put]
func (h *UserHandler) AdminSetScore(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidUserID.Error())
		return
	}
	var req AdminSetScoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误: "+err.Error())
		return
	}
	currentScore, err := h.userSvc.GetScoreById(uint(targetID))
	if err != nil {
		response.InternalError(c, "查询当前积分失败")
		return
	}
	var newScore int
	switch req.Operation {
	case "set":
		newScore = req.Score
	case "add":
		newScore = currentScore + req.Score
	case "subtract":
		newScore = currentScore - req.Score
	}
	if newScore < 0 {
		response.BadRequest(c, "积分不能为负数")
		return
	}
	if newScore > 999999 {
		response.BadRequest(c, "积分超出最大限制")
		return
	}
	viewerID, _ := c.Get("user_id")
	viewerUint, _ := viewerID.(uint)
	err = h.userSvc.SetScoreById(uint(targetID), newScore)
	if err != nil {
		response.InternalError(c, "设置积分失败: "+err.Error())
		return
	}
	response.Success(c, AdminSetScoreResponse{
		UserID:     targetID,
		OldScore:   currentScore,
		NewScore:   newScore,
		Change:     newScore - currentScore,
		Operation:  req.Operation,
		OperatorID: viewerUint,
		Reason:     req.Reason,
		Timestamp:  time.Now().Unix(),
	})
}

// AdminGetUserScore 获取用户积分
// @Summary 获取用户积分
// @Description 获取指定用户积分，不传id则获取所有用户积分列表
// @Tags 管理员后台
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query int false "用户ID"
// @Success 200 {object} vo.BasicResponse 
// @Failure 400 {object} vo.BasicResponse
// @Failure 401 {object} vo.BasicResponse
// @Failure 403 {object} vo.BasicResponse
// @Failure 500 {object} vo.BasicResponse
// @Router /admin/users/score [get]
func (h *UserHandler) AdminGetUserScore(c *gin.Context) {
	targetID := c.Query("id")
	if targetID == "" {
		users, err := h.userSvc.GetAllUsersWithScore()
		if err != nil {
			response.InternalError(c, "查询用户积分失败")
			return
		}
		response.Success(c, users)
		return
	}
	id, err := strconv.ParseUint(targetID, 10, 64)
	if err != nil {
		response.BadRequest(c, apperrors.ErrInvalidUserID.Error())
		return
	}
	score, err := h.userSvc.GetScoreById(uint(id))
	if err != nil {
		response.InternalError(c, apperrors.ErrFailedToQueryScore.Error())
		return
	}
	response.Success(c, gin.H{
		"user_id": id,
		"score":   score,
	})
}
