package board

import (
	"strconv"
	"time"

	boardService "tiny-forum/internal/service/board"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// BanUser 禁言用户
// @Summary 禁言用户
// @Description 在指定板块禁言用户，需要版主或管理员权限
// @Tags 禁言管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param body body object true "禁言信息"
// @Param body.user_id body int true "用户ID" example(10086)
// @Param body.reason body string true "禁言原因" example("发布违规内容")
// @Param body.expires_at body string false "过期时间（RFC3339格式，空表示永久）" example("2024-12-31T23:59:59Z")
// @Success 200 {object} common.BasicResponse  "禁言成功"
// @Failure 400 {object} common.BasicResponse"请求参数错误或板块ID无效"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（需要版主或管理员权限）"
// @Failure 404 {object} common.BasicResponse"板块不存在"
// @Router /boards/{id}/bans [post]
func (h *BoardHandler) BanUser(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}

	var body struct {
		UserID    uint   `json:"user_id"   binding:"required"`
		Reason    string `json:"reason"    binding:"required"`
		ExpiresAt string `json:"expires_at"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	banInput := boardService.BanUserInput{
		UserID:  body.UserID,
		BoardID: uint(boardID),
		Reason:  body.Reason,
	}
	if body.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, body.ExpiresAt)
		if err != nil {
			response.BadRequest(c, "无效的过期时间格式（需 RFC3339）")
			return
		}
		banInput.ExpiresAt = &t
	}

	operatorID := c.GetUint("user_id")
	if err := h.boardSvc.BanUser(banInput, operatorID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "禁言成功"})
}

// UnbanUser 解除禁言
// @Summary 解除禁言
// @Description 解除用户在指定板块的禁言，需要版主或管理员权限
// @Tags 禁言管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param user_id path int true "用户ID"
// @Success 200 {object} common.BasicResponse "解除禁言成功"
// @Failure 400 {object} common.BasicResponse"无效的板块ID或用户ID"
// @Failure 401 {object} common.BasicResponse"未授权"
// @Failure 403 {object} common.BasicResponse"无权限（需要版主或管理员权限）"
// @Failure 404 {object} common.BasicResponse"禁言记录不存在"
// @Router /boards/{id}/bans/{user_id} [delete]
func (h *BoardHandler) UnbanUser(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	if err := h.boardSvc.UnbanUser(uint(userID), uint(boardID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "解除禁言成功"})
}
