package user

import (
	"tiny-forum/internal/repository/token" // 假设 token 包路径

	"gorm.io/gorm"
)

type UserRepository struct {
	db        *gorm.DB
	tokenRepo token.TokenRepository
}

func NewUserRepository(db *gorm.DB, tokenRepo token.TokenRepository) *UserRepository {
	return &UserRepository{db: db, tokenRepo: tokenRepo}
}
