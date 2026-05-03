package auth

import (
	"tiny-forum/internal/infra/config"
	authService "tiny-forum/internal/service/auth"
)

type AuthHandler struct {
	// userSvc *userService.UserService
	authSvc authService.AuthService
	// jwtMgr  jwt.JWTManager
	cfg *config.Config
}

func NewAuthHandler(authSvc authService.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, cfg: cfg}
}
