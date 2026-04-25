// Package ratelimit 提供基于 Redis 的滑动窗口限流。
// 依赖：github.com/redis/go-redis/v9
package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Action 被限流的操作类型
type Action string

const (
	ActionCreatePost    Action = "create_post"
	ActionCreateComment Action = "create_comment"
	ActionSendReport    Action = "send_report"
	ActionUpdateProfile Action = "update_profile"
)

// RiskLevel 对应 model.RiskLevel，此处复制避免循环依赖
type RiskLevel string

const (
	RiskNormal   RiskLevel = "normal"   // 正常用户
	RiskObserve  RiskLevel = "observe"  // 观察用户（可疑）
	RiskRestrict RiskLevel = "restrict" // 受限用户（高危）
)

// Quota 单位时间内的配额
type Quota struct {
	Limit  int
	Window time.Duration
}

// quotaTable 各风险等级 × 操作 的配额表
// key: RiskLevel → Action → Quota
var quotaTable = map[RiskLevel]map[Action]Quota{
	// 正常用户每小时可发 20 个帖子、60 条评论、10 次举报、5 次修改资料
	RiskNormal: {
		ActionCreatePost:    {Limit: 20, Window: time.Hour},
		ActionCreateComment: {Limit: 60, Window: time.Hour},
		ActionSendReport:    {Limit: 10, Window: time.Hour},
		ActionUpdateProfile: {Limit: 5, Window: time.Hour},
	},
	// 观察用户每小时只能发 5 个帖子、20 条评论、5 次举报、3 次修改资料
	RiskObserve: {
		ActionCreatePost:    {Limit: 5, Window: time.Hour},
		ActionCreateComment: {Limit: 20, Window: time.Hour},
		ActionSendReport:    {Limit: 5, Window: time.Hour},
		ActionUpdateProfile: {Limit: 3, Window: time.Hour},
	},
	//受限用户每小时只能发 2 个帖子、5 条评论、2 次举报、1 次修改资料
	RiskRestrict: {
		ActionCreatePost:    {Limit: 2, Window: time.Hour},
		ActionCreateComment: {Limit: 5, Window: time.Hour},
		ActionSendReport:    {Limit: 2, Window: time.Hour},
		ActionUpdateProfile: {Limit: 1, Window: time.Hour},
	},
}

// Limiter 限流器
type Limiter struct {
	rdb *redis.Client
}

// Result 限流结果
type Result struct {
	Allowed bool
	Current int           // 当前窗口已用次数
	Limit   int           // 配额上限
	ResetIn time.Duration // 距离窗口重置的时间
}

func NewLimiter(rdb *redis.Client) *Limiter {
	return &Limiter{rdb: rdb}
}

// Allow 检查 userID 执行 action 是否被允许（滑动窗口算法）
// riskLevel 由调用方从 RiskService 获取
func (l *Limiter) Allow(ctx context.Context, userID uint, action Action, riskLevel RiskLevel) (Result, error) {
	quota, ok := quotaTable[riskLevel][action]
	if !ok {
		// 未配置则不限流
		return Result{Allowed: true}, nil
	}

	key := fmt.Sprintf("rl:%d:%s", userID, action)
	now := time.Now()
	windowStart := now.Add(-quota.Window)

	pipe := l.rdb.Pipeline()
	// 移除窗口外的旧记录
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart.UnixMilli()))
	// 统计当前窗口内的请求数
	countCmd := pipe.ZCard(ctx, key)
	// 设置 key 过期（防止内存泄漏）
	pipe.Expire(ctx, key, quota.Window+time.Minute)

	if _, err := pipe.Exec(ctx); err != nil {
		return Result{}, fmt.Errorf("ratelimit pipeline: %w", err)
	}

	current := int(countCmd.Val())
	if current >= quota.Limit {
		return Result{
			Allowed: false,
			Current: current,
			Limit:   quota.Limit,
			ResetIn: quota.Window,
		}, nil
	}

	// 允许通过，记录本次请求
	score := float64(now.UnixMilli())
	member := fmt.Sprintf("%d", now.UnixNano())
	if err := l.rdb.ZAdd(ctx, key, redis.Z{Score: score, Member: member}).Err(); err != nil {
		return Result{}, fmt.Errorf("ratelimit zadd: %w", err)
	}

	return Result{
		Allowed: true,
		Current: current + 1,
		Limit:   quota.Limit,
	}, nil
}

// GetQuota 查询某用户某操作的当前配额使用情况（不消耗配额）
func (l *Limiter) GetQuota(ctx context.Context, userID uint, action Action, riskLevel RiskLevel) (Result, error) {
	quota, ok := quotaTable[riskLevel][action]
	if !ok {
		return Result{Allowed: true}, nil
	}

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
