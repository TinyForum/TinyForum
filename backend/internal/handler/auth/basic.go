package auth

import (
	"tiny-forum/config"
	authService "tiny-forum/internal/service/auth"
	"tiny-forum/pkg/jwt"
)

type AuthHandler struct {
	// userSvc *userService.UserService
	authSvc authService.AuthService
	jwtMgr  jwt.JWTManager
	cfg     *config.Config
}

func NewAuthHandler(authSvc authService.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}
