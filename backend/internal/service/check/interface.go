package check

import (
	"tiny-forum/internal/infra/sensitive"
	"tiny-forum/internal/model/po"
	riskrepo "tiny-forum/internal/repository/risk"
)

type ContentCheckService interface {
	CheckPostContent(title, content string) CheckResult
	CheckText(text string) CheckResult
	CreateAuditTaskForPost(postID uint, triggerType string, hitWords []string) error
	CreateAuditTaskForComment(commentID uint, triggerType string, hitWords []string) error
	HandleReportAggregate(targetType po.AuditTargetType, targetID uint) (triggered bool, err error)
	GetListPendingTasks(limit, offset int) ([]po.ContentAuditTask, int64, error)
	ResolveTask(taskID uint, approved bool, reviewerID uint, note string) error
}

// ContentCheckService 内容安全检测服务
type contentCheckService struct {
	repo   riskrepo.RiskRepository
	filter sensitive.Filter
}

func NewContentCheckService(repo riskrepo.RiskRepository, filter sensitive.Filter) ContentCheckService {
	return &contentCheckService{repo: repo, filter: filter}
}
