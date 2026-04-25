package ratelimit

import (
	"context"
	"fmt"
	"tiny-forum/config"

	"github.com/redis/go-redis/v9"
)

type RateLimiter interface {
	Allow(ctx context.Context, userID uint, action Action, riskLevel RiskLevel) (Result, error)
	GetQuota(ctx context.Context, userID uint, action Action, riskLevel RiskLevel) (Result, error)
}

// NewLimiter 创建限流器，接收配置构建配额表
func NewLimiter(rdb *redis.Client, cfg config.RateLimitConfig) (*Limiter, error) {
	quotas, err := buildQuotaTable(cfg.RiskLevels)
	if err != nil {
		return nil, fmt.Errorf("build quota table: %w", err)
	}
	script := redis.NewScript(luaScript)
	return &Limiter{
		rdb:         rdb,
		allowScript: script,
		quotas:      quotas,
	}, nil
}
