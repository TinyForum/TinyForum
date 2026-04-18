package auth

import (
	"os"
	authService "tiny-forum/internal/service/auth"
	userService "tiny-forum/internal/service/user"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	// userSvc *userService.UserService
	authSvc authService.AuthService
}

func NewAuthHandler(authSvc authService.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Register godoc
// @Summary 用户注册
// @Tags 验证管理
// @Accept json
// @Produce json
// @Param body body user.RegisterInput true "注册信息"
// @Success 200 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var input userService.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.authSvc.Register(ctx, input)
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
// @Param body body user.LoginInput true "登录信息"
// @Success 200 {object} response.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	var input userService.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.authSvc.Login(ctx, input)
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
