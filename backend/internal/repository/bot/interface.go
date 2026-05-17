package bot

import (
	"context"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, bot *do.Bot) error
	Update(ctx context.Context, id uint, updates map[string]interface{}) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*do.Bot, error)
	ListByUser(ctx context.Context, creatorID uint, offset, limit int) ([]*do.Bot, int64, error)
	List(ctx context.Context, offset, limit int) ([]*do.Bot, int64, error)
	ListActive(ctx context.Context) ([]*do.Bot, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{db: db}
}
