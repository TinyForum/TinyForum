package risk

import (
	"time"
	"tiny-forum/internal/model/po"

	"gorm.io/gorm"
)

type RiskRepository interface {
	CreateAuditLog(log *po.AuditLog) error
	ListAuditLogs(targetType string, targetID uint, limit int) ([]po.AuditLog, error)
	CreateAuditTask(task *po.ContentAuditTask) error
	ListPendingTasks(limit, offset int) ([]po.ContentAuditTask, int64, error)
	UpdateTaskStatus(taskID uint, status po.ModerationStatus, reviewerID uint, note string) error
	CountPendingByTarget(targetType po.AuditTargetType, targetID uint) (int64, error)
	AddRiskRecord(record *po.UserRiskRecord) error
	CountActiveRiskEvents(userID uint) (int64, error)

	// IP相关
	CountActiveRiskEventsByIP(ip string) (int, error)
	AddIPRiskRecord(record *po.IPRiskRecord) error
	IsIPBlocked(ip string) (bool, error)
}
type riskRepository struct {
	db *gorm.DB
}

func NewRiskRepository(db *gorm.DB) RiskRepository {
	return &riskRepository{db: db}
}

// ========================
// AuditLog
// ========================

func (r *riskRepository) CreateAuditLog(log *po.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *riskRepository) ListAuditLogs(targetType string, targetID uint, limit int) ([]po.AuditLog, error) {
	var logs []po.AuditLog
	q := r.db.Preload("Operator").Order("created_at DESC").Limit(limit)
	if targetType != "" {
		q = q.Where("target_type = ? AND target_id = ?", targetType, targetID)
	}
	return logs, q.Find(&logs).Error
}

// ========================
// ContentAuditTask
// ========================

func (r *riskRepository) CreateAuditTask(task *po.ContentAuditTask) error {
	// 幂等：同一目标已有 pending 任务则不重复创建
	var existing po.ContentAuditTask
	err := r.db.Where("target_type = ? AND target_id = ? AND status = ?",
		task.TargetType, task.TargetID, po.ModerationStatusPending).First(&existing).Error
	if err == nil {
		return nil // 已存在，跳过
	}
	return r.db.Create(task).Error
}

func (r *riskRepository) ListPendingTasks(limit, offset int) ([]po.ContentAuditTask, int64, error) {
	var tasks []po.ContentAuditTask
	var total int64
	r.db.Model(&po.ContentAuditTask{}).Where("status = ?", po.ModerationStatusPending).Count(&total)
	err := r.db.Where("status = ?", po.ModerationStatusPending).
		Preload("Reviewer").
		Order("created_at ASC").
		Limit(limit).Offset(offset).
		Find(&tasks).Error
	return tasks, total, err
}

func (r *riskRepository) UpdateTaskStatus(taskID uint, status po.ModerationStatus, reviewerID uint, note string) error {
	now := time.Now()
	return r.db.Model(&po.ContentAuditTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":      status,
		"reviewer_id": reviewerID,
		"review_note": note,
		"reviewed_at": &now,
	}).Error
}

func (r *riskRepository) CountPendingByTarget(targetType po.AuditTargetType, targetID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.Report{}).
		Where("target_type = ? AND target_id = ? AND status = ?", targetType, targetID, po.ReportPending).
		Count(&count).Error
	return count, err
}

// ========================
// UserRiskRecord
// ========================

func (r *riskRepository) AddRiskRecord(record *po.UserRiskRecord) error {
	return r.db.Create(record).Error
}

// CountActiveRiskEvents 统计用户未过期的风险事件数，用于计算风险等级
func (r *riskRepository) CountActiveRiskEvents(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&po.UserRiskRecord{}).
		Where("user_id = ? AND expire_at > ?", userID, time.Now()).
		Count(&count).Error
	return count, err
}
