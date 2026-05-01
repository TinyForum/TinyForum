package token

import (
	"context"
	"time"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

// Repository 令牌数据访问接口
type TokenRepository interface {
	Create(ctx context.Context, userID uint, token string, expiresAt time.Time, userAgent, ip string) (*do.RefreshToken, error)
	FindByToken(ctx context.Context, token string) (*do.RefreshToken, error)
	FindByJTI(ctx context.Context, jti string) (*do.RefreshToken, error)
	FindByUserID(ctx context.Context, userID uint) ([]*do.RefreshToken, error)
	Revoke(ctx context.Context, tokenID uint) error
	RevokeAllByUserID(ctx context.Context, userID uint) error
	RevokeByJTI(ctx context.Context, jti string) error
	CleanExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uint) error
	DeleteByToken(ctx context.Context, token string) error
	DeleteByJTI(ctx context.Context, jti string) error
	DeleteExpired(ctx context.Context) (int64, error)
	DeleteByUserIDExcept(ctx context.Context, userID uint, keepToken string) error
	ListByUserID(ctx context.Context, userID uint) ([]do.RefreshToken, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	DeleteByUserIDWithTx(ctx context.Context, tx *gorm.DB, userID uint) error
	GetRefreshTokenTTL(ctx context.Context, jti string) (time.Duration, error)
	// SaveResetToken 保存密码重置 token（覆盖同用户旧 token，实现单请求生效）
	SaveResetToken(ctx context.Context, userID uint, token string, expiration time.Duration) error
	// FindUserByResetToken 通过 reset token 查找对应 user ID（自动过滤已过期）
	// FIX #41: 使用独立存储，与 user 表解耦
	FindUserByResetToken(ctx context.Context, token string) (uint, error)

	// DeleteResetToken 使用后立即删除 reset token（单次使用保证）
	// FIX #41: token 用完即废，防止重放攻击
	DeleteResetToken(ctx context.Context, token string) error

	// ---- JWT Token Blacklist ----

	// RevokeToken 将 JWT 的 JTI 加入黑名单（用于注销）
	// FIX #50/#51: 注销后 token 立即失效，防止重放攻击
	// expiresAt 用于设置黑名单记录的 TTL，节省存储空间
	RevokeToken(ctx context.Context, jti string, expiresAt time.Time) error

	// IsTokenRevoked 检查 JTI 是否已被吊销（在认证中间件中调用）
	// FIX #50: 中间件调用此方法拦截已注销的 token
	IsTokenRevoked(ctx context.Context, jti string) (bool, error)

	// ---- Session Management ----

	// RevokeAllUserTokens 吊销某用户的所有 refresh token（改密/删号后全端下线）
	// FIX #31: 修改密码后强制其他设备重新登录
	// FIX #48: 重置密码后清除所有旧会话
	RevokeAllUserTokens(ctx context.Context, userID uint) error
}
