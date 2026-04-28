package token

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type tokenRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewTokenRepository(db *gorm.DB, redis *redis.Client) TokenRepository {
	return &tokenRepository{db: db, redis: redis}
}
