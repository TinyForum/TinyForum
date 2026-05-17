package reports

import (
	"context"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type ReportsRepository interface {
	List(ctx context.Context, listReportBO *common.PageQuery[bo.ListReportBO]) ([]do.Report, int64, error)
}
type reportsRepository struct {
	db *gorm.DB
}

func NewRiskRepository(db *gorm.DB) ReportsRepository {
	return &reportsRepository{db: db}
}
