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
	titleResult := s.checker.Check(title)
	contentResult := s.checker.Check(content)
	allMatches := append(titleResult.Matches, contentResult.Matches...)
	hitWords := extractWords(allMatches)
	higher := mergeByAction(titleResult, contentResult)

	return vo.CheckResult{
		Passed:   higher.Action != sensitive.ActionBlock,
		Level:    higher.Level,
		Action:   higher.Action,
		HitWords: hitWords,
		Replaced: higher.Masked,
	}
}

func mergeByAction(a, b *sensitive.CheckResult) *sensitive.CheckResult {
	// 优先级：block > review > replace > shadow > pass
	order := map[sensitive.Action]int{
		sensitive.ActionBlock:   4,
		sensitive.ActionReview:  3,
		sensitive.ActionReplace: 2,
		sensitive.ActionShadow:  1,
		sensitive.ActionPass:    0,
	}
	if order[a.Action] >= order[b.Action] {
		return a
	}
	return b
}

func extractWords(matches []*sensitive.MatchResult) []string {
	words := make([]string, 0, len(matches))
	for _, m := range matches {
		if m.Word != "" {
			words = append(words, m.Word)
		}
	}
	return words
}

// CheckText 检测单段文本（评论、简介等）
func (s *contentCheckService) CheckText(text string) vo.CheckResult {
	r := s.checker.Check(text)
	return vo.CheckResult{
		Passed:   r.Action != sensitive.ActionBlock,
		Level:    r.Level,
		Action:   r.Action,
		HitWords: extractWords(r.Matches),
		Replaced: r.Masked,
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
