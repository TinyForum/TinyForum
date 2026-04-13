package handler

import (
	"os"
	"tiny-forum/internal/service"
	"tiny-forum/pkg/response"
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
// @Tags 验证管理
// @Accept json
// @Produce json
// @Param body body service.RegisterInput true "注册信息"
// @Success 200 {object} response.Response
// @Router /auth/register [post]
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
// @Tags 验证管理
// @Accept json
// @Produce json
// @Param body body service.LoginInput true "登录信息"
// @Success 200 {object} response.Response
// @Router /auth/login [post]
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

	isProduction := os.Getenv("APP_ENV") == "production"

	c.SetCookie(
		"tiny_forum_token",
		result.Token,
		3600*24*7,
		"/",
		"",
		isProduction, // 生产环境强制 HTTPS
		true,         // HttpOnly，JS 无法读取
	)

	// 响应体只返回用户信息，不暴露 token
	response.Success(c, gin.H{
		"user": result.User,
	})
}

// Logout godoc
// @Summary 用户登出
// @Tags 验证管理
// @Produce json
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	isProduction := os.Getenv("APP_ENV") == "production"

	// MaxAge=-1 立即删除 Cookie
	c.SetCookie(
		"tiny_forum_token",
		"",
		-1,
		"/",
		"",
		isProduction,
		true,
	)
	response.Success(c, nil)
}

// Me godoc
// @Summary 获取当前用户信息
// @Tags 验证管理
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userSvc.GetProfile(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}