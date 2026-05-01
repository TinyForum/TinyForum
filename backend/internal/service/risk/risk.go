package risk

import (
	"context"
	"fmt"
	"time"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/model/do"
)

// GetUserRiskLevel 计算用户当前风险等级
// 规则（优先级从高到低）：
//  1. IsBlocked → blocked（调用方直接拦截，无需走此函数）
//  2. 活跃风险事件 >= 3 → restrict
//  3. score < 50 或 注册不足7天 → observe
//  4. 其他 → normal
func (s *riskService) GetUserRiskLevel(user *do.User) (do.RiskLevel, error) {
	if user.IsBlocked {
		return do.RiskLevelBlocked, nil
	}

	activeEvents, err := s.repo.CountActiveRiskEvents(user.ID)
	if err != nil {
		return do.RiskLevelNormal, err
	}
	if activeEvents >= 3 {
		return do.RiskLevelRestrict, nil
	}

	isNewUser := time.Since(user.CreatedAt) < 7*24*time.Hour
	isLowScore := user.Score < 50
	if isNewUser || isLowScore {
		return do.RiskLevelObserve, nil
	}

	return do.RiskLevelNormal, nil
}

// toRatelimitLevel 将 do.RiskLevel 转换为 ratelimit 包的类型（避免循环依赖）
func toRatelimitLevel(level do.RiskLevel) ratelimit.RiskLevel {
	switch level {
	case do.RiskLevelRestrict:
		return ratelimit.RiskRestrict
	case do.RiskLevelObserve:
		return ratelimit.RiskObserve
	default:
		return ratelimit.RiskNormal
	}
}

// CheckRateLimit 检查用户操作频率是否超限
func (s *riskService) CheckRateLimit(ctx context.Context, user *do.User, action ratelimit.Action) (ratelimit.Result, error) {
	level, err := s.GetUserRiskLevel(user)
	if err != nil {
		return ratelimit.Result{Allowed: true}, nil // 降级放行
	}
	return s.limiter.Allow(ctx, fmt.Sprint(user.ID), action, toRatelimitLevel(level))
}

// RecordRiskEvent 记录一次风险事件（举报成立、命中敏感词等）
// ttl: 该事件计入风险分的有效期
func (s *riskService) RecordRiskEvent(userID uint, eventType, detail string, ttl time.Duration) error {
	record := &do.UserRiskRecord{
		UserID:      userID,
		EventType:   eventType,
		EventDetail: detail,
		ExpireAt:    time.Now().Add(ttl),
	}
	return s.repo.AddRiskRecord(record)
}

// WriteAuditLog 写入操作审计日志
func (s *riskService) WriteAuditLog(operatorID uint, action do.AuditActionType,
	targetType string, targetID uint, before, after, reason, ip string) error {
	log := &do.AuditLog{
		OperatorID: operatorID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Before:     before,
		After:      after,
		Reason:     reason,
		IP:         ip,
	}
	return s.repo.CreateAuditLog(log)
}
