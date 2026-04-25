package user

import (
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"
	jwtpkg "tiny-forum/pkg/jwt"
	"tiny-forum/pkg/validator"
)

type UserService struct {
	repo        userRepo.UserRepository
	jwtMgr      *jwtpkg.Manager
	notifSvc    *notification.NotificationService // 注意：NotificationService 定义在别的包，需正确导入
	roleChecker *validator.RoleChangeChecker      // 改为指针类型
	// roleChange  validator.RoleChangeRequest
}

func NewUserService(
	repo userRepo.UserRepository,
	jwtMgr *jwtpkg.Manager,
	notifSvc *notification.NotificationService,
	// roleChange validator.RoleChangeRequest,
) *UserService {
	roleValidator := validator.NewRoleValidator()
	roleChecker := validator.NewRoleChangeChecker(roleValidator)
	return &UserService{
		repo:     repo,
		jwtMgr:   jwtMgr,
		notifSvc: notifSvc,
		// roleChecker: validator.RoleChangeChecker{},
		// roleChange:  validator.RoleChangeRequest{},
		roleChecker: roleChecker,
	}
}

// avatarURL 生成默认头像URL
func avatarURL(username string) string {
	return "https://api.dicebear.com/8.x/lorelei/svg?seed=" + username
}
