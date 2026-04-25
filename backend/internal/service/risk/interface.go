package risk

import (
	"context"
	"time"
	"tiny-forum/internal/model"
	riskrepo "tiny-forum/internal/repository/risk"
	"tiny-forum/pkg/ratelimit"
)

type RiskService interface {
	GetAuditLogs(targetType string, targetID uint, limit int) ([]model.AuditLog, error)
	GetUserRiskLevel(user *model.User) (model.RiskLevel, error)
	CheckRateLimit(ctx context.Context, user *model.User, action ratelimit.Action) (ratelimit.Result, error)
	RecordRiskEvent(userID uint, eventType, detail string, ttl time.Duration) error
	WriteAuditLog(operatorID uint, action model.AuditActionType,
		targetType string, targetID uint, before, after, reason, ip string) error
}

// RiskService 风控核心服务
type riskService struct {
	repo    riskrepo.RiskRepository
	limiter ratelimit.RateLimiter
}

func NewRiskService(repo riskrepo.RiskRepository, limiter ratelimit.RateLimiter) RiskService {
	return &riskService{repo: repo, limiter: limiter}
}
