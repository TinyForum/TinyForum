package risk

import (
	"context"
	"time"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/model/do"
	riskrepo "tiny-forum/internal/repository/risk"
)

type RiskService interface {
	GetAuditLogs(targetType string, targetID uint, limit int) ([]do.AuditLog, error)
	GetUserRiskLevel(user *do.User) (do.RiskLevel, error)
	CheckRateLimit(ctx context.Context, user *do.User, action ratelimit.Action) (ratelimit.Result, error)
	RecordRiskEvent(userID uint, eventType, detail string, ttl time.Duration) error
	WriteAuditLog(operatorID uint, action do.AuditActionType,
		targetType string, targetID uint, before, after, reason, ip string) error
	GetAnonymousRiskLevel(ip string) (do.RiskLevel, error)
	// ip
	CheckRateLimitByIP(ctx context.Context, ip string, action ratelimit.Action) (ratelimit.Result, error)
	RecordRiskEventByIP(ip, eventType, detail string, ttl time.Duration) error
	WriteAuditLogByIP(ip string, action do.AuditActionType, targetType string,
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
