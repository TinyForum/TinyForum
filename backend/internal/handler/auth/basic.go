package auth

import (
	authService "tiny-forum/internal/service/auth"
)

type AuthHandler struct {
	// userSvc *userService.UserService
	authSvc authService.AuthService
	// jwtMgr  jwt.JWTManager
	// cfg     *config.Config
}

func NewAuthHandler(authSvc authService.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}
