package ratelimit

import (
	"context"
	_ "embed"
	"fmt"
	"time"

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
	rdb         *redis.Client
	allowScript *redis.Script // Lua 脚本缓存
}

// Result 限流结果
type Result struct {
	Allowed bool
	Current int           // 当前窗口已用次数
	Limit   int           // 配额上限
	ResetIn time.Duration // 距离窗口重置的时间
}

func NewLimiter(rdb *redis.Client) *Limiter {
	script := redis.NewScript(luaScript)
	return &Limiter{
		rdb:         rdb,
		allowScript: script,
	}
}

// Allow 检查 userID 执行 action 是否被允许（滑动窗口算法）
// riskLevel 由调用方从 RiskService 获取
// Allow 原子版：使用 Lua 脚本实现滑动窗口限流
func (l *Limiter) Allow(ctx context.Context, userID uint, action Action, riskLevel RiskLevel) (Result, error) {
	quota, ok := quotaTable[riskLevel][action]
	if !ok {
		// 未配置则不限流
		return Result{Allowed: true}, nil
	}

	key := fmt.Sprintf("rl:%d:%s", userID, action)
	now := time.Now()
	windowMs := quota.Window.Milliseconds()
	nowMs := now.UnixMilli()
	member := fmt.Sprintf("%d", now.UnixNano()) // 纳秒级唯一标识

	// 执行 Lua 脚本
	res, err := l.allowScript.Run(ctx, l.rdb, []string{key},
		quota.Limit, windowMs, nowMs, member).Slice()

	if err != nil {
		return Result{}, fmt.Errorf("ratelimit lua script: %w", err)
	}

	// 脚本返回: [allowed, current, limit, resetMs]

	if len(res) != 4 {
		return Result{}, fmt.Errorf("unexpected script result length: %d", len(res))
	}

	// 类型转换辅助
	toInt := func(v interface{}) (int, error) {
		switch t := v.(type) {
		case int64:
			return int(t), nil
		case int:
			return t, nil
		default:
			return 0, fmt.Errorf("cannot convert %T to int", v)
		}
	}
	toInt64 := func(v interface{}) (int64, error) {
		switch t := v.(type) {
		case int64:
			return t, nil
		case int:
			return int64(t), nil
		default:
			return 0, fmt.Errorf("cannot convert %T to int64", v)
		}
	}

	allowed, err := toInt(res[0])
	if err != nil {
		return Result{}, err
	}
	current, err := toInt(res[1])
	if err != nil {
		return Result{}, err
	}
	limit, err := toInt(res[2])
	if err != nil {
		return Result{}, err
	}
	resetMs, err := toInt64(res[3])
	if err != nil {
		return Result{}, err
	}
	fmt.Printf("Result: allowed=%v, current=%d, limit=%d, resetIn=%v\n",
		res[0].(int64) == 1, res[1].(int64), res[2].(int64), time.Duration(res[3].(int64))*time.Millisecond)
	result := Result{
		Allowed: allowed == 1,
		Current: current,
		Limit:   limit,
	}
	if !result.Allowed {
		result.ResetIn = time.Duration(resetMs) * time.Millisecond
	}
	return result, nil
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
