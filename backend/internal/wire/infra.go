// internal/wire/infra.go
//
// 变更说明：
//   - Infra 新增 Enforcer *casbin.Enforcer 字段
//   - InitInfra 新增 db *gorm.DB 参数（enforcer 需要数据库连接）
//   - 其余逻辑与原版完全一致

package wire

import (
	"context"
	"log"
	"time"

	"tiny-forum/internal/infra/casbinx"
	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/infra/sensitive"

	"github.com/casbin/casbin/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Infra 封装基础设施（面向接口）
type Infra struct {
	RedisClient     *redis.Client
	RateLimiter     ratelimit.RateLimiter
	SensitiveFilter sensitive.Filter
	Enforcer        *casbin.Enforcer // Casbin RBAC enforcer
}

// InitInfra 初始化基础设施。
//
// db 参数用于 Casbin 的 GORM adapter，会在数据库中自动创建 casbin_rule 表。
func InitInfra(cfg *config.Config, db *gorm.DB) (*Infra, error) {
	log.Println("初始化基础设施...")
	// 1. Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	// 2. 限流器
	rateLimiter, err := ratelimit.NewLimiter(redisClient, cfg.RiskControl.RateLimit)
	if err != nil {
		return nil, err
	}

	// 3. 敏感词过滤器
	ollamaCfg := &sensitive.OllamaConfig{
		BaseURL: cfg.Basic.Ollama.BaseURL,
		Model:   cfg.Basic.Ollama.Model,
		Timeout: time.Duration(cfg.Basic.Ollama.Timeout) * time.Second,
	}
	sensitiveFilter := sensitive.NewFilter(ollamaCfg)
	res, err := sensitiveFilter.LoadDictDir("./dicts")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("词库加载完成：block %d 个文件，review %d 个文件，共 %d 词条，失败 %d 个",
		len(res.BlockFiles), len(res.ReviewFiles), res.TotalWords, len(res.Errors))

	// 4. Casbin enforcer（策略持久化到 casbin_rule 表）
	enforcer, err := casbinx.NewEnforcer(db, "config/rbac_model.conf")
	if err != nil {
		return nil, err
	}

	return &Infra{
		RedisClient:     redisClient,
		RateLimiter:     rateLimiter,
		SensitiveFilter: sensitiveFilter,
		Enforcer:        enforcer,
	}, nil
}
