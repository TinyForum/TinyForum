// package auth

// import (
// 	"tiny-forum/config"
// 	userRepo "tiny-forum/internal/repository/user"
// 	"tiny-forum/internal/service/email"
// 	"tiny-forum/internal/service/notification"
// 	jwtpkg "tiny-forum/pkg/jwt"
// 	"tiny-forum/pkg/validator"
// )

// type AuthService struct {
// 	repo        *userRepo.UserRepository
// 	jwtMgr      *jwtpkg.Manager
// 	notifSvc    *notification.NotificationService
// 	roleChecker validator.RoleChangeChecker
// 	emailSvc    email.Service
// 	cfg         *config.Config
// }

// func NewAuthService(
//
//	repo *userRepo.UserRepository,
//	jwtMgr *jwtpkg.Manager,
//	notifSvc *notification.NotificationService,
//	emailSvc email.Service,
//	cfg *config.Config,
//
//	) *AuthService {
//		return &AuthService{
//			repo:        repo,
//			jwtMgr:      jwtMgr,
//			notifSvc:    notifSvc,
//			roleChecker: validator.RoleChangeChecker{},
//			emailSvc:    emailSvc,
//			cfg:         cfg,
//		}
//	}
package auth
