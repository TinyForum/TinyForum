package wire

import (
	"log"
	"time"

	"tiny-forum/config"
	"tiny-forum/pkg/ratelimit"
	"tiny-forum/pkg/sensitive"

	"github.com/redis/go-redis/v9"
)

// Infra 封装 Redis、限流器、敏感词过滤器等基础设施
type Infra struct {
	RedisClient     *redis.Client
	RateLimiter     *ratelimit.Limiter
	SensitiveFilter sensitive.Filter
}

// InitInfra 初始化基础设施
func InitInfra(cfg *config.Config) (*Infra, error) {
	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Private.Redis.GetAddr(),
		Password: cfg.Private.Redis.Password,
		DB:       cfg.Private.Redis.DB,
	})
	rateLimiter := ratelimit.NewLimiter(rdb)

	// 敏感词检测
	ollamaCfg := &sensitive.OllamaConfig{
		BaseURL: cfg.Basic.Ollama.BaseURL,
		Model:   cfg.Basic.Ollama.Model,
		Timeout: 15 * time.Second,
	}
	sensitiveFilter := sensitive.NewFilter(ollamaCfg)
	res, err := sensitiveFilter.LoadDictDir("./dicts")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("词库加载完成：block %d 个文件，review %d 个文件，共 %d 词条，失败 %d 个",
		len(res.BlockFiles), len(res.ReviewFiles), res.TotalWords, len(res.Errors))

	return &Infra{
		RedisClient:     rdb,
		RateLimiter:     rateLimiter,
		SensitiveFilter: sensitiveFilter,
	}, nil
}
