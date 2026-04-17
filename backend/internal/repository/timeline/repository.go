package timeline

import (
	"gorm.io/gorm"
)

type TimelineRepository struct {
	db *gorm.DB
}

func NewTimelineRepository(db *gorm.DB) *TimelineRepository {
	return &TimelineRepository{db: db}
}
