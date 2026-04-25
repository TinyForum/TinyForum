package stats

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type statsRepository struct {
	db *gorm.DB
}

type StatsRepository interface {
	GetHotArticlesByDateRange(
		ctx context.Context,
		startDate, endDate time.Time,
		limit int,
	) ([]*HotArticleRow, error)
	GetHotBoardsByDateRange(
		ctx context.Context,
		startDate, endDate time.Time,
		limit int,
	) ([]*HotBoardRow, error)
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepository{db: db}
}
