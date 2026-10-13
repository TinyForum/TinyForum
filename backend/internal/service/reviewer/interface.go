package reviewer

import (
	"context"

	"tiny-forum/internal/model/common"
	"tiny-forum/internal/model/vo"
)

// ReviewerService 审核员业务接口
type ReviewerService interface {
	// ListPendingTasks 分页获取待审核任务
	ListPendingTasks(ctx context.Context, page, pageSize int) (*common.PageResult[vo.AuditTaskVO], error)
	// ListReviewedTasks 分页获取已处理的审核任务（审核历史）
	ListReviewedTasks(ctx context.Context, page, pageSize int) (*common.PageResult[vo.AuditTaskVO], error)
	// ReviewTask 处理审核任务：通过或拒绝，并联动更新目标内容状态
	ReviewTask(ctx context.Context, taskID, reviewerID uint, approved bool, note string) error
}
