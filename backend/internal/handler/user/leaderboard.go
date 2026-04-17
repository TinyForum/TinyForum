package user

import (
	"strings"
	"tiny-forum/internal/model"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Leaderboard 用户排行榜 [great]
// @Summary 获取用户排行榜
// @Tags 用户管理
// @Produce json
// @Param limit query int false "数量" default(20)
// @Param fields query string false "需要返回的字段，逗号分隔" default(id,username,avatar,score)
// @Router /users/leaderboard [get]
func (h *UserHandler) Leaderboard(c *gin.Context) {
	var req LeaderboardRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 字段白名单校验
	if req.Fields != "" {
		allowedMap := make(map[string]bool)
		for _, f := range model.UserPublicFields {
			allowedMap[f] = true
		}
		fields := strings.Split(req.Fields, ",")
		for _, f := range fields {
			f = strings.TrimSpace(f)
			if f == "" {
				continue
			}
			if !allowedMap[f] {
				response.BadRequest(c, "无效的字段名: "+f)
				return
			}
		}
	}

	items, err := h.userSvc.GetLeaderboard(c.Request.Context(), req.Limit, req.Fields)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	resp := LeaderboardResponse{
		Items: make([]LeaderboardItemResponse, len(items)),
	}
	for i, item := range items {
		resp.Items[i] = LeaderboardItemResponse{
			ID:       item.ID,
			Username: item.Username,
			Avatar:   item.Avatar,
			Score:    item.Score,
			Rank:     item.Rank,
		}
	}
	response.Success(c, resp)
}
