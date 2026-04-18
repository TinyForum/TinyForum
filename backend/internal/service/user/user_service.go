package user

import (
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"
	jwtpkg "tiny-forum/pkg/jwt"
	"tiny-forum/pkg/validator"
)

type UserService struct {
	repo        *userRepo.UserRepository
	jwtMgr      *jwtpkg.Manager
	notifSvc    *notification.NotificationService // 注意：NotificationService 定义在别的包，需正确导入
	roleChecker validator.RoleChangeChecker
}

func NewUserService(
	repo *userRepo.UserRepository,
	jwtMgr *jwtpkg.Manager,
	notifSvc *notification.NotificationService,
) *UserService {
	return &UserService{
		repo:        repo,
		jwtMgr:      jwtMgr,
		notifSvc:    notifSvc,
		roleChecker: validator.RoleChangeChecker{},
	}
}

// avatarURL 生成默认头像URL
func avatarURL(username string) string {
	return "https://api.dicebear.com/8.x/lorelei/svg?seed=" + username
}
