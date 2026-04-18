// internal/auth/repository/repository.go
package auth

import (
	"context"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *authRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *authRepository) FindByResetToken(ctx context.Context, token string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("reset_password_token = ?", token).First(&user).Error
	return &user, err
}

func (r *authRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *authRepository) Save(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
