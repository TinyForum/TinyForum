// internal/auth/interface.go
package auth

import (
	"context"
	"time"
	"tiny-forum/config"
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository/auth"
	"tiny-forum/internal/repository/token"
	"tiny-forum/internal/repository/transaction"
	"tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/email"
	"tiny-forum/internal/service/notification"
	userSvc "tiny-forum/internal/service/user"
	jwtpkg "tiny-forum/pkg/jwt"

	"github.com/redis/go-redis/v9"
)

type DeleteAccountInput struct {
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}
type DeletionStatus struct {
	IsDeleted     bool       `json:"is_deleted"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	CanRestore    bool       `json:"can_restore"`
	RemainingDays int        `json:"remaining_days,omitempty"`
}

// Service 定义业务逻辑接口
type AuthService interface {
	ForgotPassword(ctx context.Context, email, locale string) error
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
	ValidateResetToken(ctx context.Context, token string) (bool, error)
	Login(ctx context.Context, input userSvc.LoginInput) (*AuthResult, error)
	Register(ctx context.Context, input userSvc.RegisterInput) (*userSvc.AuthResult, error)
	ChangePassword(userID uint, oldPassword, newPassword string) (string, error)
	DeleteAccount(ctx context.Context, userID uint, input DeleteAccountInput) error
	CancelDeletion(ctx context.Context, userID uint) error
	ConfirmDeletion(ctx context.Context, userID uint) error
	GetDeletionStatus(ctx context.Context, userID uint) (*DeletionStatus, error)
	RevokeToken(ctx context.Context, jti string) error
}
type authService struct {
	userRepo  user.UserRepository
	notifSvc  notification.NotificationService
	jwtMgr    *jwtpkg.JWTManager
	authRepo  auth.AuthRepository
	emailSvc  email.EmailService
	cfg       *config.Config
	tokenRepo token.TokenRepository
	txManager transaction.TransactionManager
	redis     *redis.Client
}

// Repository 定义数据访问接口
type Repository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByResetToken(ctx context.Context, token string) (*model.User, error)
	FindByID(ctx context.Context, id uint) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Create(ctx context.Context, user *model.User) error
}

func NewAuthService(
	authRepo auth.AuthRepository,
	userRepo user.UserRepository,
	jwtMgr *jwtpkg.JWTManager,
	notifSvc notification.NotificationService,
	emailSvc email.EmailService,
	cfg *config.Config,
	tokenRepo token.TokenRepository,
	txManager transaction.TransactionManager,
	redis *redis.Client,
) AuthService {

	return &authService{
		authRepo:  authRepo,
		userRepo:  userRepo,
		notifSvc:  notifSvc,
		jwtMgr:    jwtMgr,
		emailSvc:  emailSvc,
		cfg:       cfg,
		tokenRepo: tokenRepo,
		txManager: txManager,
		redis:     redis,
	}
}
