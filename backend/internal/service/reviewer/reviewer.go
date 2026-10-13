package reviewer

import (
	"context"
	"errors"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
	commentRepo "tiny-forum/internal/repository/comment"
	articleService "tiny-forum/internal/service/article"
	checkService "tiny-forum/internal/service/check"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

type reviewerService struct {
	checkSvc    checkService.ContentCheckService
	articleSvc  articleService.ArticleService
	commentRepo commentRepo.CommentRepository
}

// NewReviewerService 创建审核员服务
func NewReviewerService(
	checkSvc checkService.ContentCheckService,
	articleSvc articleService.ArticleService,
	commentRepo commentRepo.CommentRepository,
) ReviewerService {
	return &reviewerService{
		checkSvc:    checkSvc,
		articleSvc:  articleSvc,
		commentRepo: commentRepo,
	}
}

// ListPendingTasks 分页获取待审核任务
func (s *reviewerService) ListPendingTasks(ctx context.Context, page, pageSize int) (*common.PageResult[vo.AuditTaskVO], error) {
	offset := (page - 1) * pageSize
	tasks, total, err := s.checkSvc.GetListPendingTasks(pageSize, offset)
	if err != nil {
		return nil, err
	}

	list := make([]vo.AuditTaskVO, 0, len(tasks))
	for _, t := range tasks {
		list = append(list, toAuditTaskVO(t))
	}

	return &common.PageResult[vo.AuditTaskVO]{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  int64(offset+len(list)) < total,
	}, nil
}

// ListReviewedTasks 分页获取已处理的审核任务（审核历史）
func (s *reviewerService) ListReviewedTasks(ctx context.Context, page, pageSize int) (*common.PageResult[vo.AuditTaskVO], error) {
	offset := (page - 1) * pageSize
	tasks, total, err := s.checkSvc.GetResolvedTasks(pageSize, offset)
	if err != nil {
		return nil, err
	}

	list := make([]vo.AuditTaskVO, 0, len(tasks))
	for _, t := range tasks {
		list = append(list, toAuditTaskVO(t))
	}

	return &common.PageResult[vo.AuditTaskVO]{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  int64(offset+len(list)) < total,
	}, nil
}

// toAuditTaskVO 将审核任务实体转换为视图对象
func toAuditTaskVO(t do.ContentAuditTask) vo.AuditTaskVO {
	task := vo.AuditTaskVO{
		ID:          t.ID,
		CreatedAt:   t.CreatedAt,
		TargetType:  string(t.TargetType),
		TargetID:    t.TargetID,
		TriggerType: string(t.TriggerType),
		TriggerMeta: t.TriggerMeta,
		Status:      string(t.Status),
		ReviewerID:  t.ReviewerID,
		ReviewNote:  t.ReviewNote,
		ReviewedAt:  t.ReviewedAt,
	}
	if t.Reviewer != nil {
		task.ReviewerName = t.Reviewer.Username
	}
	return task
}

// ReviewTask 处理审核任务并联动更新目标内容状态
func (s *reviewerService) ReviewTask(ctx context.Context, taskID, reviewerID uint, approved bool, note string) error {
	task, err := s.checkSvc.GetTaskByID(taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound.WithDetail("审核任务不存在")
		}
		return err
	}
	if task.Status != do.ModerationStatusPending {
		return apperrors.ErrInvalidRequest.WithDetail("该审核任务已被处理")
	}

	// 更新任务状态
	if err := s.checkSvc.ResolveTask(taskID, approved, reviewerID, note); err != nil {
		return err
	}

	// 联动更新目标内容状态
	switch task.TargetType {
	case do.AuditTargetPost:
		status := do.ModerationStatusRejected
		if approved {
			status = do.ModerationStatusApproved
		}
		return s.articleSvc.AdminSetReviewPost(task.TargetID, status)
	case do.AuditTargetComment:
		status := do.ReplyStatusHidden
		if approved {
			status = do.ReplyStatusVisible
		}
		return s.commentRepo.UpdateReplyStatusByCommentID(ctx, task.TargetID, status)
	default:
		// 用户资料审核暂无联动动作
		return nil
	}
}
