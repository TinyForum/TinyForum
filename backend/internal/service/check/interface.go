package check

import (
	"tiny-forum/internal/infra/sensitive"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
	riskrepo "tiny-forum/internal/repository/risk"
)

type ContentCheckService interface {
	CheckPostContent(title, content string) vo.CheckResult
	CheckText(text string) vo.CheckResult
	CreateAuditTaskForPost(postID uint, triggerType string, hitWords []string) error
	CreateAuditTaskForComment(commentID uint, triggerType string, hitWords []string) error
	HandleReportAggregate(targetType do.AuditTargetType, targetID uint) (triggered bool, err error)
	GetListPendingTasks(limit, offset int) ([]do.ContentAuditTask, int64, error)
	GetResolvedTasks(limit, offset int) ([]do.ContentAuditTask, int64, error)
	GetTaskByID(taskID uint) (*do.ContentAuditTask, error)
	ResolveTask(taskID uint, approved bool, reviewerID uint, note string) error
	// CheckAndCreateCommentTask 检测评论内容，命中规则则创建审核任务
	CheckAndCreateCommentTask(commentID uint, content string) error
}

// ContentCheckService 内容安全检测服务
type contentCheckService struct {
	repo    riskrepo.RiskRepository
	checker *sensitive.Checker
}

func NewContentCheckService(repo riskrepo.RiskRepository, checker *sensitive.Checker) ContentCheckService {
	return &contentCheckService{repo: repo, checker: checker}
}
