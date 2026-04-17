package user

import (
	"tiny-forum/internal/repository"
	"tiny-forum/internal/service/notification"
	jwtpkg "tiny-forum/pkg/jwt"
)

type UserService struct {
	repo        *repository.UserRepository
	jwtMgr      *jwtpkg.Manager
	notifSvc    *notification.NotificationService // 注意：NotificationService 定义在别的包，需正确导入
	roleChecker RoleChangeChecker
}

func NewUserService(
	repo *repository.UserRepository,
	jwtMgr *jwtpkg.Manager,
	notifSvc *notification.NotificationService,
) *UserService {
	return &UserService{
		repo:        repo,
		jwtMgr:      jwtMgr,
		notifSvc:    notifSvc,
		roleChecker: RoleChangeChecker{},
	}
}

// avatarURL 生成默认头像URL
func avatarURL(username string) string {
	return "https://api.dicebear.com/8.x/lorelei/svg?seed=" + username
}
