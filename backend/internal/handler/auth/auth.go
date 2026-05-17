package auth

import (
	"log"
	"os"
	userService "tiny-forum/internal/service/user"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary 用户登录
// @Tags 验证管理
// @Accept json
// @Produce json
// @Param body body user.LoginInput true "登录信息"
// @Success 200 {object} vo.UserPrivateVO
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {

	ctx := c.Request.Context()
	var input userService.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.HandleError(c, err)
		return
	}

	result, err := h.authSvc.Login(ctx, input)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	isProduction := os.Getenv("APP_ENV") == "production"

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
	log.Printf("JWT secret used for signing: %s", h.cfg.Private.JWT.Secret)
	log.Printf("JWT secret length: %d", len(h.cfg.Private.JWT.Secret))

	c.Header("Set-Cookie", cookieValue)

	response.Success(c, result.User)
}

// Logout godoc
// @Summary 用户登出
// @Tags 验证管理
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} common.BasicResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	isProduction := os.Getenv("APP_ENV") == "production"

	// 从 Cookie 中取出 token，通知 service 层使其失效
	if token, err := c.Cookie("tiny_forum_token"); err == nil && token != "" {
		// 将 token 加入黑名单
		_ = h.authSvc.RevokeToken(ctx, token)
	}

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
