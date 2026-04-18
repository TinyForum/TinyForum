package token

import (
	"context"
	"time"
	"tiny-forum/internal/model"
)

// Repository 令牌数据访问接口
type TokenRepository interface {
	Create(ctx context.Context, userID uint, token string, expiresAt time.Time, userAgent, ip string) (*model.RefreshToken, error)
	FindByToken(ctx context.Context, token string) (*model.RefreshToken, error)
	FindByJTI(ctx context.Context, jti string) (*model.RefreshToken, error)
	FindByUserID(ctx context.Context, userID uint) ([]*model.RefreshToken, error)
	Revoke(ctx context.Context, tokenID uint) error
	RevokeAllByUserID(ctx context.Context, userID uint) error
	RevokeByJTI(ctx context.Context, jti string) error
	CleanExpired(ctx context.Context) error
	DeleteByUserID(ctx context.Context, userID uint) error
	DeleteByToken(ctx context.Context, token string) error
	DeleteByJTI(ctx context.Context, jti string) error
	DeleteExpired(ctx context.Context) (int64, error)
	DeleteByUserIDExcept(ctx context.Context, userID uint, keepToken string) error
	ListByUserID(ctx context.Context, userID uint) ([]model.RefreshToken, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	SaveResetToken(ctx context.Context, userID uint, token string, expiration time.Duration) error
}
