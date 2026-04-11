package handler

import (
	"bbs-forum/internal/service"
	"bbs-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userSvc *service.UserService
}

func NewAuthHandler(userSvc *service.UserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

// Register godoc
// @Summary 用户注册
// @Tags auth
// @Accept json
// @Produce json
// @Param body body service.RegisterInput true "注册信息"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.userSvc.Register(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// Login godoc
// @Summary 用户登录
// @Tags auth
// @Accept json
// @Produce json
// @Param body body service.LoginInput true "登录信息"
// @Success 200 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.userSvc.Login(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	c.SetCookie(
		"bbs_token",
		result.Token,
		3600*24*7, // 7 days
		"/",
		"",
		false,
		true, // HttpOnly
	)
	response.Success(c, result)
}

// Me godoc
// @Summary 获取当前用户信息
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userSvc.GetProfile(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}
