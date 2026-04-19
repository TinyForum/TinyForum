package risk

import (
	"time"
	"tiny-forum/internal/model"

	"gorm.io/gorm"
)

type RiskRepository struct {
	db *gorm.DB
}

func NewRiskRepository(db *gorm.DB) *RiskRepository {
	return &RiskRepository{db: db}
}

// ========================
// AuditLog
// ========================

func (r *RiskRepository) CreateAuditLog(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *RiskRepository) ListAuditLogs(targetType string, targetID uint, limit int) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	q := r.db.Preload("Operator").Order("created_at DESC").Limit(limit)
	if targetType != "" {
		q = q.Where("target_type = ? AND target_id = ?", targetType, targetID)
	}
	return logs, q.Find(&logs).Error
}

// ========================
// ContentAuditTask
// ========================

func (r *RiskRepository) CreateAuditTask(task *model.ContentAuditTask) error {
	// 幂等：同一目标已有 pending 任务则不重复创建
	var existing model.ContentAuditTask
	err := r.db.Where("target_type = ? AND target_id = ? AND status = ?",
		task.TargetType, task.TargetID, model.ModerationStatusPending).First(&existing).Error
	if err == nil {
		return nil // 已存在，跳过
	}
	return r.db.Create(task).Error
}

func (r *RiskRepository) ListPendingTasks(limit, offset int) ([]model.ContentAuditTask, int64, error) {
	var tasks []model.ContentAuditTask
	var total int64
	r.db.Model(&model.ContentAuditTask{}).Where("status = ?", model.ModerationStatusPending).Count(&total)
	err := r.db.Where("status = ?", model.ModerationStatusPending).
		Preload("Reviewer").
		Order("created_at ASC").
		Limit(limit).Offset(offset).
		Find(&tasks).Error
	return tasks, total, err
}

func (r *RiskRepository) UpdateTaskStatus(taskID uint, status model.ModerationStatus, reviewerID uint, note string) error {
	now := time.Now()
	return r.db.Model(&model.ContentAuditTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":      status,
		"reviewer_id": reviewerID,
		"review_note": note,
		"reviewed_at": &now,
	}).Error
}

func (r *RiskRepository) CountPendingByTarget(targetType model.AuditTargetType, targetID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Report{}).
		Where("target_type = ? AND target_id = ? AND status = ?", targetType, targetID, model.ReportPending).
		Count(&count).Error
	return count, err
}

// ========================
// UserRiskRecord
// ========================

func (r *RiskRepository) AddRiskRecord(record *model.UserRiskRecord) error {
	return r.db.Create(record).Error
}

// CountActiveRiskEvents 统计用户未过期的风险事件数，用于计算风险等级
func (r *RiskRepository) CountActiveRiskEvents(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.UserRiskRecord{}).
		Where("user_id = ? AND expire_at > ?", userID, time.Now()).
		Count(&count).Error
	return count, err
}
