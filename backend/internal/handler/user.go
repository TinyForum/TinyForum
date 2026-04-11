package handler

import (
	"strconv"

	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}

	viewerID, _ := c.Get("user_id")
	var viewerUint uint
	if v, ok := viewerID.(uint); ok {
		viewerUint = v
	}

	profile, err := h.userSvc.GetUserProfile(uint(targetID), viewerUint)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input service.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.userSvc.UpdateProfile(userID, input); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	user, _ := h.userSvc.GetProfile(userID)
	response.Success(c, user)
}

func (h *UserHandler) Follow(c *gin.Context) {
	followerID := c.GetUint("user_id")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	if err := h.userSvc.Follow(followerID, uint(targetID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "关注成功"})
}

func (h *UserHandler) Unfollow(c *gin.Context) {
	followerID := c.GetUint("user_id")
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	if err := h.userSvc.Unfollow(followerID, uint(targetID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "已取消关注"})
}

func (h *UserHandler) Leaderboard(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	users, err := h.userSvc.GetLeaderboard(limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, users)
}

// Admin handlers
func (h *UserHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	users, total, err := h.userSvc.List(page, pageSize, keyword)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, users, total, page, pageSize)
}

func (h *UserHandler) AdminSetActive(c *gin.Context) {
	targetID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	var body struct {
		Active bool `json:"active"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.userSvc.SetActive(uint(targetID), body.Active); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "操作成功"})
}
