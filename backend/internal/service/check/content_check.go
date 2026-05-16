package check

import (
	"encoding/json"
	"fmt"
	"tiny-forum/internal/infra/sensitive"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
	apperrors "tiny-forum/pkg/errors"
)

// CheckPostContent 检测帖子内容（title + content）
func (s *contentCheckService) CheckPostContent(title, content string) vo.CheckResult {
	// title 和 content 分别检测，取最高等级
	titleResult := s.filter.Check(title)
	contentResult := s.filter.Check(content)

	higher := titleResult
	if contentResult.Level > titleResult.Level {
		higher = contentResult
	}

	return vo.CheckResult{
		Passed:   higher.Level != sensitive.LevelBlock,
		Level:    higher.Level,
		HitWords: append(titleResult.HitWords, contentResult.HitWords...),
		Replaced: higher.Text,
	}
}

// CheckText 检测单段文本（评论、简介等）
func (s *contentCheckService) CheckText(text string) vo.CheckResult {
	r := s.filter.Check(text)
	return vo.CheckResult{
		Passed:   r.Level != sensitive.LevelBlock,
		Level:    r.Level,
		HitWords: r.HitWords,
		Replaced: r.Text,
	}
}

// CreateAuditTaskForPost 为帖子创建审核任务
func (s *contentCheckService) CreateAuditTaskForPost(postID uint, triggerType string, hitWords []string) error {
	meta, _ := json.Marshal(map[string]interface{}{
		"hit_words": hitWords,
	})
	auditTriggerType, err := do.ParseAuditTriggerType(triggerType)
	if err != nil {
		return apperrors.ErrValidation
	}
	task := &do.ContentAuditTask{
		TargetType:  do.AuditTargetPost,
		TargetID:    postID,
		TriggerType: auditTriggerType,
		TriggerMeta: string(meta),
		Status:      do.ModerationStatusPending,
	}
	return s.repo.CreateAuditTask(task)
}

// CreateAuditTaskForComment 为评论创建审核任务
func (s *contentCheckService) CreateAuditTaskForComment(commentID uint, triggerType string, hitWords []string) error {
	meta, _ := json.Marshal(map[string]interface{}{
		"hit_words": hitWords,
	})
	auditTriggerType, err := do.ParseAuditTriggerType(triggerType)
	if err != nil {
		return apperrors.ErrValidation
	}
	task := &do.ContentAuditTask{
		TargetType:  do.AuditTargetComment,
		TargetID:    commentID,
		TriggerType: auditTriggerType,
		TriggerMeta: string(meta),
		Status:      do.ModerationStatusPending,
	}
	return s.repo.CreateAuditTask(task)
}

// HandleReportAggregate 处理举报聚合逻辑
// 当某内容的 pending 举报数达到阈值时，自动创建审核任务并将内容状态改为 pending
// 返回是否触发了聚合
func (s *contentCheckService) HandleReportAggregate(
	targetType do.AuditTargetType, targetID uint,
) (triggered bool, err error) {
	count, err := s.repo.CountPendingByTarget(targetType, targetID)
	if err != nil {
		return false, err
	}
	if count < do.ReportAggregateThreshold {
		return false, nil
	}

	// 创建审核任务
	meta, _ := json.Marshal(map[string]interface{}{
		"report_count": count,
	})
	task := &do.ContentAuditTask{
		TargetType:  targetType,
		TargetID:    targetID,
		TriggerType: "report_aggregate",
		TriggerMeta: string(meta),
		Status:      do.ModerationStatusPending,
	}
	if err = s.repo.CreateAuditTask(task); err != nil {
		return false, fmt.Errorf("create audit task: %w", err)
	}
	return true, nil
}

// GetListPendingTasks 获取待审核任务列表
func (s *contentCheckService) GetListPendingTasks(limit, offset int) ([]do.ContentAuditTask, int64, error) {
	return s.repo.ListPendingTasks(limit, offset)
}

// ResolveTask 处理审核任务
func (s *contentCheckService) ResolveTask(taskID uint, approved bool, reviewerID uint, note string) error {
	status := do.ModerationStatusApproved
	if !approved {
		status = do.ModerationStatusRejected
	}
	return s.repo.UpdateTaskStatus(taskID, status, reviewerID, note)
}
