// internal/auth/repository/interface.go
package auth

import (
	"context"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByResetToken(ctx context.Context, token string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Save(ctx context.Context, user *model.User) error
	SoftDelete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	HardDelete(ctx context.Context, id uint) error
	HardDeleteWithTx(ctx context.Context, tx *gorm.DB, id uint) error
	GetDeletedUser(ctx context.Context, id uint) (*model.User, error)
	GetUserWithDeleted(ctx context.Context, id uint) (*model.User, error)
}
type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{db: db}
}
