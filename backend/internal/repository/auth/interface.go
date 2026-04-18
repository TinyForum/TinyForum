// internal/auth/repository/interface.go
package auth

import (
	"context"
	"tiny-forum/internal/model"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByResetToken(ctx context.Context, token string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Save(ctx context.Context, user *model.User) error
}
