package wire

import (
	"context"
	"log"
	"time"

	"tiny-forum/config"
	"tiny-forum/pkg/ratelimit"
	"tiny-forum/pkg/sensitive/filter"

	"github.com/redis/go-redis/v9"
)

// Infra 封装基础设施（面向接口）
type Infra struct {
	RedisClient     *redis.Client
	RateLimiter     ratelimit.RateLimiter // 使用接口而非具体实现
	SensitiveFilter filter.Filter
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
	ollamaCfg := &filter.OllamaConfig{
		BaseURL: cfg.Basic.Ollama.BaseURL,
		Model:   cfg.Basic.Ollama.Model,
		Timeout: 15 * time.Second,
	}
	sensitiveFilter := filter.NewFilter(ollamaCfg)
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
