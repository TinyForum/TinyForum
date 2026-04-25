package check

import (
	"encoding/json"
	"fmt"
	"tiny-forum/internal/model"
	"tiny-forum/pkg/sensitive"
)

// CheckResult 内容检测结果
type CheckResult struct {
	Passed   bool            // false 表示直接拦截
	Level    sensitive.Level // 命中等级
	HitWords []string        // 命中词
	Replaced string          // 替换后的内容
}

// CheckPostContent 检测帖子内容（title + content）
func (s *contentCheckService) CheckPostContent(title, content string) CheckResult {
	// title 和 content 分别检测，取最高等级
	titleResult := s.filter.Check(title)
	contentResult := s.filter.Check(content)

	higher := titleResult
	if contentResult.Level > titleResult.Level {
		higher = contentResult
	}

	return CheckResult{
		Passed:   higher.Level != sensitive.LevelBlock,
		Level:    higher.Level,
		HitWords: append(titleResult.HitWords, contentResult.HitWords...),
		Replaced: higher.Text,
	}
}

// CheckText 检测单段文本（评论、简介等）
func (s *contentCheckService) CheckText(text string) CheckResult {
	r := s.filter.Check(text)
	return CheckResult{
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
	task := &model.ContentAuditTask{
		TargetType:  model.AuditTargetPost,
		TargetID:    postID,
		TriggerType: triggerType,
		TriggerMeta: string(meta),
		Status:      model.ModerationStatusPending,
	}
	return s.repo.CreateAuditTask(task)
}

// CreateAuditTaskForComment 为评论创建审核任务
func (s *contentCheckService) CreateAuditTaskForComment(commentID uint, triggerType string, hitWords []string) error {
	meta, _ := json.Marshal(map[string]interface{}{
		"hit_words": hitWords,
	})
	task := &model.ContentAuditTask{
		TargetType:  model.AuditTargetComment,
		TargetID:    commentID,
		TriggerType: triggerType,
		TriggerMeta: string(meta),
		Status:      model.ModerationStatusPending,
	}
	return s.repo.CreateAuditTask(task)
}

// HandleReportAggregate 处理举报聚合逻辑
// 当某内容的 pending 举报数达到阈值时，自动创建审核任务并将内容状态改为 pending
// 返回是否触发了聚合
func (s *contentCheckService) HandleReportAggregate(
	targetType model.AuditTargetType, targetID uint,
) (triggered bool, err error) {
	count, err := s.repo.CountPendingByTarget(targetType, targetID)
	if err != nil {
		return false, err
	}
	if count < model.ReportAggregateThreshold {
		return false, nil
	}

	// 创建审核任务
	meta, _ := json.Marshal(map[string]interface{}{
		"report_count": count,
	})
	task := &model.ContentAuditTask{
		TargetType:  targetType,
		TargetID:    targetID,
		TriggerType: "report_aggregate",
		TriggerMeta: string(meta),
		Status:      model.ModerationStatusPending,
	}
	if err = s.repo.CreateAuditTask(task); err != nil {
		return false, fmt.Errorf("create audit task: %w", err)
	}
	return true, nil
}

// GetListPendingTasks 获取待审核任务列表
func (s *contentCheckService) GetListPendingTasks(limit, offset int) ([]model.ContentAuditTask, int64, error) {
	return s.repo.ListPendingTasks(limit, offset)
}

// ResolveTask 处理审核任务
func (s *contentCheckService) ResolveTask(taskID uint, approved bool, reviewerID uint, note string) error {
	status := model.ModerationStatusApproved
	if !approved {
		status = model.ModerationStatusRejected
	}
	return s.repo.UpdateTaskStatus(taskID, status, reviewerID, note)
}
