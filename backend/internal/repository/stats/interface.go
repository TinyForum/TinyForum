package stats

import (
	"context"
	"time"
	"tiny-forum/internal/model/vo"

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
	) ([]*vo.HotArticleRowVO, error)
	GetHotBoardsByDateRange(
		ctx context.Context,
		startDate, endDate time.Time,
		limit int,
	) ([]*vo.HotBoardRowVO, error)
}

func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepository{db: db}
}
