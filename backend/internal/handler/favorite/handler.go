package favorite

import (
	"strconv"

	"tiny-forum/internal/model/request"
	"tiny-forum/internal/service/favorite"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// FavoriteHandler 收藏处理器
type FavoriteHandler struct {
	service favorite.FavoriteService
}

// NewFavoriteHandler 创建收藏处理器
func NewFavoriteHandler(service favorite.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

// AddFavorite 收藏内容
// @Summary 收藏内容
// @Description 收藏指定的作品（帖子/文章）
// @Tags 收藏
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body request.AddFavoriteRequest true "收藏信息"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /users/me/favorites [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	var req request.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.AddFavorite(c.Request.Context(), userID, req.TargetID, req.TargetType, req.GroupID); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"favorited": true})
}

// RemoveFavorite 取消收藏
// @Summary 取消收藏
// @Tags 收藏
// @Produce json
// @Security ApiKeyAuth
// @Param target_id path int true "收藏对象ID"
// @Success 200 {object} common.BasicResponse
// @Failure 404 {object} common.BasicResponse
// @Router /users/me/favorites/{target_id} [delete]
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("target_id"), 10, 64)
	if err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}
	userID := c.GetUint("user_id")
	if err := h.service.RemoveFavorite(c.Request.Context(), userID, uint(targetID), c.Query("target_type")); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, gin.H{"favorited": false})
}

// ListMyFavorites 我的收藏列表
// @Summary 我的收藏列表
// @Tags 收藏
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} common.BasicResponse
// @Router /users/me/favorites [get]
func (h *FavoriteHandler) ListMyFavorites(c *gin.Context) {
	var req request.ListFavoritesRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}
	userID := c.GetUint("user_id")
	result, err := h.service.ListMyFavorites(c.Request.Context(), req, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}
