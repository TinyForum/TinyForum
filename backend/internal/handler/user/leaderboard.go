package user

import (
	"tiny-forum/internal/dto"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// user_handler.go

// LeaderboarSimple 精简排行榜信息
// @Summary 精简排行榜信息
// @Tags 用户排行榜
// @Security BearerAuth
// @Param limit query int false "限制返回数量"
// @Success 200 {array} dto.SimpleLeaderboardItem
// @Router /users/leaderboard/simple [get]
func (h *UserHandler) LeaderboardSimple(c *gin.Context) {
	var req dto.LeaderboardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 调用 service 获取原始数据（返回 model.User 切片或自定义结构）
	users, err := h.userSvc.GetSimpleLeaderboardData(c.Request.Context(), req.Limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	items := make([]dto.SimpleLeaderboardItem, len(users))
	for i, u := range users {
		items[i] = dto.SimpleLeaderboardItem{
			ID:       u.ID,
			Username: u.Username,
			Score:    u.Score,
			Rank:     i + 1,
		}
	}
	response.Success(c, items)
}

// LeaderboardDetail 详细排行榜信息
// @Summary 详细排行榜信息
// @Tags 用户排行榜
// @Security BearerAuth
// @Param limit query int false "限制返回数量"
// @Success 200 {array} dto.DetailLeaderboardItem
// @Router /users/leaderboard/detail [get]
func (h *UserHandler) LeaderboardDetail(c *gin.Context) {
	var req dto.LeaderboardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	users, err := h.userSvc.GetDetailLeaderboardData(c.Request.Context(), req.Limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	items := make([]dto.LeaderboardUserDetail, len(users))
	for i, u := range users {
		items[i] = dto.LeaderboardUserDetail{
			ID:       u.ID,
			Username: u.Username,
			Avatar:   u.Avatar,
			Email:    u.Email,
			Role:     u.Role,
			Score:    u.Score,
			Rank:     i + 1,
		}
	}
	response.Success(c, items)
}
