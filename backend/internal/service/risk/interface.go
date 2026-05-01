package risk

import (
	"context"
	"time"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/model/po"
	riskrepo "tiny-forum/internal/repository/risk"
)

type RiskService interface {
	GetAuditLogs(targetType string, targetID uint, limit int) ([]po.AuditLog, error)
	GetUserRiskLevel(user *po.User) (po.RiskLevel, error)
	CheckRateLimit(ctx context.Context, user *po.User, action ratelimit.Action) (ratelimit.Result, error)
	RecordRiskEvent(userID uint, eventType, detail string, ttl time.Duration) error
	WriteAuditLog(operatorID uint, action po.AuditActionType,
		targetType string, targetID uint, before, after, reason, ip string) error
	GetAnonymousRiskLevel(ip string) (po.RiskLevel, error)
	// ip
	CheckRateLimitByIP(ctx context.Context, ip string, action ratelimit.Action) (ratelimit.Result, error)
	RecordRiskEventByIP(ip, eventType, detail string, ttl time.Duration) error
	WriteAuditLogByIP(ip string, action po.AuditActionType, targetType string,
		targetID uint, before, after, reason string) error
}

// RiskService 风控核心服务
type riskService struct {
	repo    riskrepo.RiskRepository
	limiter ratelimit.RateLimiter
}

func NewRiskService(repo riskrepo.RiskRepository, limiter ratelimit.RateLimiter) RiskService {
	return &riskService{repo: repo, limiter: limiter}
}
