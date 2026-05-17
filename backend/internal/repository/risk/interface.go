package risk

import (
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

type RiskRepository interface {
	CreateAuditLog(log *do.AuditLog) error
	ListAuditLogs(targetType string, targetID uint, limit int) ([]do.AuditLog, error)
	CreateAuditTask(task *do.ContentAuditTask) error
	ListPendingTasks(limit, offset int) ([]do.ContentAuditTask, int64, error)
	UpdateTaskStatus(taskID uint, status do.ModerationStatus, reviewerID uint, note string) error
	CountPendingByTarget(targetType do.AuditTargetType, targetID uint) (int64, error)
	AddRiskRecord(record *do.UserRiskRecord) error
	CountActiveRiskEvents(userID uint) (int64, error)

	// IP相关
	CountActiveRiskEventsByIP(ip string) (int, error)
	AddIPRiskRecord(record *do.IPRiskRecord) error
	IsIPBlocked(ip string) (bool, error)
}
type riskRepository struct {
	db *gorm.DB
}

func NewRiskRepository(db *gorm.DB) RiskRepository {
	return &riskRepository{db: db}
}
