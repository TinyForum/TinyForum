package wire

import (
	"context"
	"log"
	"time"

	"tiny-forum/config"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/infra/sensitive"

	"github.com/redis/go-redis/v9"
)

// Infra 封装基础设施（面向接口）
type Infra struct {
	RedisClient     *redis.Client
	RateLimiter     ratelimit.RateLimiter
	SensitiveFilter sensitive.Filter
}

// InitInfra 初始化基础设施
func InitInfra(cfg *config.Config) (*Infra, error) {
	// 1. Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Private.Redis.GetAddr(),
		Password: cfg.Private.Redis.Password,
		DB:       cfg.Private.Redis.DB,
	})
	// 可选：测试连接
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	// 2. 限流器（传入风控配置）
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
		log.Fatal(err) // 词库加载失败是严重问题，可直接终止
	}
	log.Printf("词库加载完成：block %d 个文件，review %d 个文件，共 %d 词条，失败 %d 个",
		len(res.BlockFiles), len(res.ReviewFiles), res.TotalWords, len(res.Errors))

	return &Infra{
		RedisClient:     redisClient,
		RateLimiter:     rateLimiter,
		SensitiveFilter: sensitiveFilter,
	}, nil
}
