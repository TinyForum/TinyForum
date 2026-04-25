package user

import (
	authService "tiny-forum/internal/service/auth"
	notiService "tiny-forum/internal/service/notification"
	userService "tiny-forum/internal/service/user"
)

type UserHandler struct {
	userSvc  userService.UserService
	notifSvc notiService.NotificationService
	authSvc  authService.AuthService
}

func NewUserHandler(userSvc userService.UserService, notifSvc notiService.NotificationService, authSvc authService.AuthService) *UserHandler {
	return &UserHandler{userSvc: userSvc, notifSvc: notifSvc, authSvc: authSvc}
}
