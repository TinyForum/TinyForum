package risk

import (
	"time"
	"tiny-forum/internal/model/do"
)

// ========================
// AuditLog
// ========================

func (r *riskRepository) CreateAuditLog(log *do.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *riskRepository) ListAuditLogs(targetType string, targetID uint, limit int) ([]do.AuditLog, error) {
	var logs []do.AuditLog
	q := r.db.Preload("Operator").Order("created_at DESC").Limit(limit)
	if targetType != "" {
		q = q.Where("target_type = ? AND target_id = ?", targetType, targetID)
	}
	return logs, q.Find(&logs).Error
}

// ========================
// ContentAuditTask
// ========================

func (r *riskRepository) CreateAuditTask(task *do.ContentAuditTask) error {
	// 幂等：同一目标已有 pending 任务则不重复创建
	var existing do.ContentAuditTask
	err := r.db.Where("target_type = ? AND target_id = ? AND status = ?",
		task.TargetType, task.TargetID, do.ModerationStatusPending).First(&existing).Error
	if err == nil {
		return nil // 已存在，跳过
	}
	return r.db.Create(task).Error
}

func (r *riskRepository) ListPendingTasks(limit, offset int) ([]do.ContentAuditTask, int64, error) {
	var tasks []do.ContentAuditTask
	var total int64
	r.db.Model(&do.ContentAuditTask{}).Where("status = ?", do.ModerationStatusPending).Count(&total)
	err := r.db.Where("status = ?", do.ModerationStatusPending).
		Preload("Reviewer").
		Order("created_at ASC").
		Limit(limit).Offset(offset).
		Find(&tasks).Error
	return tasks, total, err
}

// ListResolvedTasks 分页获取已处理（通过/拒绝）的审核任务
func (r *riskRepository) ListResolvedTasks(limit, offset int) ([]do.ContentAuditTask, int64, error) {
	var tasks []do.ContentAuditTask
	var total int64
	resolved := []do.ModerationStatus{do.ModerationStatusApproved, do.ModerationStatusRejected}
	r.db.Model(&do.ContentAuditTask{}).Where("status IN ?", resolved).Count(&total)
	err := r.db.Where("status IN ?", resolved).
		Preload("Reviewer").
		Order("reviewed_at DESC").
		Limit(limit).Offset(offset).
		Find(&tasks).Error
	return tasks, total, err
}

// GetTaskByID 根据 ID 获取审核任务（预加载审核人）
func (r *riskRepository) GetTaskByID(taskID uint) (*do.ContentAuditTask, error) {
	var task do.ContentAuditTask
	err := r.db.Preload("Reviewer").First(&task, taskID).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *riskRepository) UpdateTaskStatus(taskID uint, status do.ModerationStatus, reviewerID uint, note string) error {
	now := time.Now()
	return r.db.Model(&do.ContentAuditTask{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"status":      status,
		"reviewer_id": reviewerID,
		"review_note": note,
		"reviewed_at": &now,
	}).Error
}

func (r *riskRepository) CountPendingByTarget(targetType do.AuditTargetType, targetID uint) (int64, error) {
	var count int64
	err := r.db.Model(&do.Report{}).
		Where("target_type = ? AND target_id = ? AND status = ?", targetType, targetID, do.ReportPending).
		Count(&count).Error
	return count, err
}

// ========================
// UserRiskRecord
// ========================

func (r *riskRepository) AddRiskRecord(record *do.UserRiskRecord) error {
	return r.db.Create(record).Error
}

// CountActiveRiskEvents 统计用户未过期的风险事件数，用于计算风险等级
func (r *riskRepository) CountActiveRiskEvents(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&do.UserRiskRecord{}).
		Where("user_id = ? AND expire_at > ?", userID, time.Now()).
		Count(&count).Error
	return count, err
}
