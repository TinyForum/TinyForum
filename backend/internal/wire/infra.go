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
	"os"
	"time"

	"tiny-forum/internal/infra/casbinx"
	"tiny-forum/internal/infra/config"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/infra/sensitive"
	"tiny-forum/pkg/logger"

	"github.com/casbin/casbin/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Infra 封装基础设施（面向接口）
type Infra struct {
	RedisClient      *redis.Client
	RateLimiter      ratelimit.RateLimiter
	sensitiveChecker *sensitive.Checker
	Enforcer         *casbin.Enforcer // Casbin RBAC enforcer
}

// InitInfra 初始化基础设施。
//
// db 参数用于 Casbin 的 GORM adapter，会在数据库中自动创建 casbin_rule 表。
func InitInfra(cfg *config.Config, db *gorm.DB) (*Infra, error) {
	logger.Infof("完整配置: %+v", cfg)
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
	logger.Infof("设置 Redis")

	// 2. 限流器
	rateLimiter, err := ratelimit.NewLimiter(redisClient, cfg.RiskControl.RateLimit)
	if err != nil {
		return nil, err
	}
	logger.Infof("设置限流器")

	// 3. 敏感词过滤器
	dictDir := "sensitive-dicts"

	if _, err = os.Stat(dictDir); os.IsNotExist(err) {
		dictDir = os.Getenv("DICT_DIR")
		if dictDir == "" {
			logger.Fatal("词典目录不存在，请设置 DICT_DIR 环境变量")
		}
	}
	aiCfg := &sensitive.AIConfig{
		Enable:   cfg.AI.Enable,
		Provider: cfg.AI.Provider,
		APIKey:   cfg.AI.Config.APIKey,
		BaseURL:  cfg.AI.Config.BaseURL,
		Model:    cfg.AI.Config.Model,
		Timeout:  time.Duration(cfg.AI.Config.Timeout) * time.Millisecond,
	}

	logger.Infof("AI 配置: %v", aiCfg)
	if aiCfg.Enable {
		logger.Infof("LLM 复判已启用，提供商: %s, 模型: %s", aiCfg.Provider, aiCfg.Model)
	} else {
		logger.Infof("LLM 复判已禁用")
	}

	sensitiveChecker, err := sensitive.NewSensitiveChecker(dictDir, aiCfg)
	if err != nil {
		logger.Fatalf("初始化审核器失败: %v", err)
	}
	logger.Infof("词典加载成功，目录: %s", dictDir)

	// 4. Casbin enforcer（策略持久化到 casbin_rule 表）
	enforcer, err := casbinx.NewEnforcer(db, "config/rbac_model.conf")
	if err != nil {
		return nil, err
	}

	return &Infra{
		RedisClient:      redisClient,
		RateLimiter:      rateLimiter,
		sensitiveChecker: sensitiveChecker,
		Enforcer:         enforcer,
	}, nil
}
