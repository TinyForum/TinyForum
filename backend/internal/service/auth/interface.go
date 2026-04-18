// internal/auth/interface.go
package auth

import (
	"context"
	"tiny-forum/config"
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository/auth"
	"tiny-forum/internal/repository/token"
	"tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/email"
	"tiny-forum/internal/service/notification"
	userSvc "tiny-forum/internal/service/user"
	jwtpkg "tiny-forum/pkg/jwt"
)

// Service 定义业务逻辑接口
type AuthService interface {
	ForgotPassword(ctx context.Context, email, locale string) error
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
	ValidateResetToken(ctx context.Context, token string) (bool, error)
	Login(ctx context.Context, input userSvc.LoginInput) (*userSvc.AuthResult, error)
	Register(ctx context.Context, input userSvc.RegisterInput) (*userSvc.AuthResult, error)
	ChangePassword(userID uint, oldPassword, newPassword string) (string, error)
}

type authService struct {
	userRepo  *user.UserRepository
	notifSvc  *notification.NotificationService
	jwtMgr    *jwtpkg.Manager
	authRepo  auth.AuthRepository
	emailSvc  email.EmailService
	cfg       *config.Config
	tokenRepo token.TokenRepository // 添加 tokenRepo
	// tokenExpiry time.Duration         // 添加这个字段
}

// Repository 定义数据访问接口
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByResetToken(ctx context.Context, token string) (*model.User, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Create(ctx context.Context, user *model.User) error
}

// EmailSender 定义邮件发送接口
type EmailSender interface {
	SendResetPasswordEmail(ctx context.Context, to, token, username, locale string) error
	SendWelcomeEmail(ctx context.Context, to, username, locale string) error
}

// TokenGenerator 定义令牌生成接口
type TokenGenerator interface {
	GenerateToken() (string, error)
}

func NewAuthService(
	userRepo *user.UserRepository,
	jwtMgr *jwtpkg.Manager,
	notifSvc *notification.NotificationService,
	emailSvc email.EmailService,
	cfg *config.Config,
	tokenRepo token.TokenRepository, // 添加参数
) AuthService {
	return &authService{
		userRepo:  userRepo,
		notifSvc:  notifSvc,
		jwtMgr:    jwtMgr,
		emailSvc:  emailSvc,
		cfg:       cfg,
		tokenRepo: tokenRepo,
	}
}
