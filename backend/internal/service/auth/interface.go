// internal/auth/interface.go
package auth

import (
	"context"
	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
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

// Service 定义业务逻辑接口
type AuthService interface {

	// create

	Register(ctx context.Context, input request.RegisterRequest) (*vo.AuthResultVO, error) // 注册

	// update

	ForgotPassword(ctx context.Context, email, ip, userAgent, locale string) error                    // 忘记密码，发送重置密码邮件
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error                           // 重置密码
	ValidateResetToken(ctx context.Context, token string) (bool, error)                               // 验证重置密码token
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) (string, error) // 修改密码
	CancelDeletion(ctx context.Context, userID uint) error                                            // 取消删除账户
	ConfirmDeletion(ctx context.Context, userID uint) error                                           // 确认删除账户
	GetUserEmailByResetToken(ctx context.Context, token string) (string, error)                       // 根据重置密码token获取用户邮箱
	ResetPasswordWithToken(ctx context.Context, token, newPassword string) error                      // 根据重置密码token重置密码

	// delete

	DeleteAccount(ctx context.Context, userID uint, input request.DeleteAccountRequest) (bool, error) // 删除账户
	RevokeToken(ctx context.Context, jti string) error                                                // 注销token（登出）

	// query

	Login(ctx context.Context, input userSvc.LoginInput) (*vo.AuthResultVO, error)  // 登录
	GetDeletionStatus(ctx context.Context, userID uint) (*vo.DeletionStatus, error) // 获取删除账户状态
	FinduUserEmailByID(userID uint) (string, error)                                 // 根据用户ID查找用户邮箱
	IsUserExist(ctx context.Context, email string) (bool, error)                    // 检查用户是否存在
	ValidateOldPassword(userID uint, newPassword string) (bool, error)              // 验证密码是否合规

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
	FindByEmail(ctx context.Context, email string) (*do.User, error)
	FindByResetToken(ctx context.Context, token string) (*do.User, error)
	FindByID(ctx context.Context, id uint) (*do.User, error)
	Update(ctx context.Context, user *do.User) error
	Create(ctx context.Context, user *do.User) error
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
