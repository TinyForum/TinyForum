package auth

import (
	"os"
	userService "tiny-forum/internal/service/user"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
// @Summary 用户注册
// @Tags 验证管理
// @Accept json
// @Produce json
// @Param body body user.RegisterInput true "注册信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	var input userService.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
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

	// FIX #57: 补充 SameSite 属性防止 CSRF
	// gin 的 SetCookie 不直接支持 SameSite，需要手动拼接 Set-Cookie header
	cookieValue := "tiny_forum_token=" + result.Token +
		"; Max-Age=604800" + // 7天
		"; Path=/" +
		// "; Domain=" + h.cfg.Basic.Server.Host +
		"; HttpOnly" +
		"; SameSite=Lax"
	if isProduction {
		cookieValue += "; Secure"
	}
	c.Header("Set-Cookie", cookieValue)

	// 响应体只返回用户信息，不暴露 token
	response.Success(c, gin.H{
		"user": result.User,
	})
}

// Logout godoc
// @Summary 用户登出
// @Tags 验证管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	isProduction := os.Getenv("APP_ENV") == "production"

	// FIX #50/#51: 服务端将 Token 加入黑名单，防止注销后 Token 重放攻击
	// 从 Cookie 中取出 token，通知 service 层使其失效
	if token, err := c.Cookie("tiny_forum_token"); err == nil && token != "" {
		// 将 token 加入黑名单（忽略错误，不影响注销流程）
		_ = h.authSvc.RevokeToken(ctx, token)
	}

	// FIX #57: 清除 Cookie 时同样补充 SameSite 属性
	cookieValue := "tiny_forum_token=" +
		"; Max-Age=-1" +
		"; Path=/" +
		"; HttpOnly" +
		"; SameSite=Strict"
	if isProduction {
		cookieValue += "; Secure"
	}
	c.Header("Set-Cookie", cookieValue)

	response.Success(c, nil)
}
