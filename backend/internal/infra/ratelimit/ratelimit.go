package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"time"
	"tiny-forum/config"

	"github.com/redis/go-redis/v9"
)

//go:embed script.lua
var luaScript string

// Action 被限流的操作类型
type Action string

const (
	ActionCreatePost    Action = "create_post"
	ActionCreateComment Action = "create_comment"
	ActionSendReport    Action = "send_report"
	ActionUpdateProfile Action = "update_profile"
)

// RiskLevel 用户风险等级
type RiskLevel string

const (
	RiskNormal   RiskLevel = "normal"
	RiskObserve  RiskLevel = "observe"
	RiskRestrict RiskLevel = "restrict"
)

// Quota 单位时间内的配额
type Quota struct {
	Limit  int
	Window time.Duration
}

// Result 限流结果
type Result struct {
	Allowed bool
	Current int           // 当前窗口已用次数
	Limit   int           // 配额上限
	ResetIn time.Duration // 距离窗口重置的时间
}

// Limiter 限流器实现
type Limiter struct {
	rdb         *redis.Client
	allowScript *redis.Script
	quotas      map[RiskLevel]map[Action]Quota // 从配置构建
}

// buildQuotaTable 将配置中的 map[string]map[string]QuotaConfig 转换为内部结构
func buildQuotaTable(riskLevels map[string]map[string]config.QuotaConfig) (map[RiskLevel]map[Action]Quota, error) {
	result := make(map[RiskLevel]map[Action]Quota)
	for riskStr, actionMap := range riskLevels {
		risk := RiskLevel(riskStr)
		result[risk] = make(map[Action]Quota)
		for actionStr, qc := range actionMap {
			action := Action(actionStr)
			window, err := time.ParseDuration(qc.Window)
			if err != nil {
				return nil, fmt.Errorf("invalid window for %s.%s: %w", riskStr, actionStr, err)
			}
			result[risk][action] = Quota{
				Limit:  qc.Limit,
				Window: window,
			}
		}
	}
	// 如果配置缺失，返回空 map，但 Allow 方法会按“不限流”处理
	return result, nil
}

// Allow 实现 RateLimiter 接口
func (l *Limiter) Allow(ctx context.Context, userID uint, action Action, riskLevel RiskLevel) (Result, error) {
	// 获取对应的配额
	if level, ok := l.quotas[riskLevel]; ok {
		if quota, ok := level[action]; ok {
			return l.checkAndRecord(ctx, userID, action, quota)
		}
	}
	// 未配置则不限流
	return Result{Allowed: true}, nil
}

// checkAndRecord 执行 Lua 脚本进行原子限流判断与记录
func (l *Limiter) checkAndRecord(ctx context.Context, userID uint, action Action, quota Quota) (Result, error) {
	key := fmt.Sprintf("rl:%d:%s", userID, action)
	now := time.Now()
	windowMs := quota.Window.Milliseconds()
	nowMs := now.UnixMilli()
	member := fmt.Sprintf("%d", now.UnixNano())

	res, err := l.allowScript.Run(ctx, l.rdb, []string{key},
		quota.Limit, windowMs, nowMs, member).Slice()
	if err != nil {
		return Result{}, fmt.Errorf("lua script error: %w", err)
	}
	if len(res) != 4 {
		return Result{}, fmt.Errorf("unexpected script result length: %d", len(res))
	}

	// 类型转换（假设返回均为 int64）
	allowed := toInt64(res[0]) == 1
	current := int(toInt64(res[1]))
	limit := int(toInt64(res[2]))
	resetMs := toInt64(res[3])

	result := Result{
		Allowed: allowed,
		Current: current,
		Limit:   limit,
	}
	if !allowed {
		result.ResetIn = time.Duration(resetMs) * time.Millisecond
	}
	return result, nil
}

// GetQuota 实现 RateLimiter 接口（不消耗）
func (l *Limiter) GetQuota(ctx context.Context, userID uint, action Action, riskLevel RiskLevel) (Result, error) {
	if level, ok := l.quotas[riskLevel]; ok {
		if quota, ok := level[action]; ok {
			key := fmt.Sprintf("rl:%d:%s", userID, action)
			windowStart := time.Now().Add(-quota.Window)
			l.rdb.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixMilli()))
			current, err := l.rdb.ZCard(ctx, key).Result()
			if err != nil {
				return Result{}, err
			}
			return Result{
				Allowed: int(current) < quota.Limit,
				Current: int(current),
				Limit:   quota.Limit,
			}, nil
		}
	}
	return Result{Allowed: true}, nil
}

// 辅助类型转换
func toInt64(v interface{}) int64 {
	switch t := v.(type) {
	case int64:
		return t
	case int:
		return int64(t)
	default:
		return 0
	}
}
