package risk

import (
	"context"
	"fmt"
	"time"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/model/po"
)

// GetUserRiskLevel 计算用户当前风险等级
// 规则（优先级从高到低）：
//  1. IsBlocked → blocked（调用方直接拦截，无需走此函数）
//  2. 活跃风险事件 >= 3 → restrict
//  3. score < 50 或 注册不足7天 → observe
//  4. 其他 → normal
func (s *riskService) GetUserRiskLevel(user *po.User) (po.RiskLevel, error) {
	if user.IsBlocked {
		return po.RiskLevelBlocked, nil
	}

	activeEvents, err := s.repo.CountActiveRiskEvents(user.ID)
	if err != nil {
		return po.RiskLevelNormal, err
	}
	if activeEvents >= 3 {
		return po.RiskLevelRestrict, nil
	}

	isNewUser := time.Since(user.CreatedAt) < 7*24*time.Hour
	isLowScore := user.Score < 50
	if isNewUser || isLowScore {
		return po.RiskLevelObserve, nil
	}

	return po.RiskLevelNormal, nil
}

// toRatelimitLevel 将 po.RiskLevel 转换为 ratelimit 包的类型（避免循环依赖）
func toRatelimitLevel(level po.RiskLevel) ratelimit.RiskLevel {
	switch level {
	case po.RiskLevelRestrict:
		return ratelimit.RiskRestrict
	case po.RiskLevelObserve:
		return ratelimit.RiskObserve
	default:
		return ratelimit.RiskNormal
	}
}

// CheckRateLimit 检查用户操作频率是否超限
func (s *riskService) CheckRateLimit(ctx context.Context, user *po.User, action ratelimit.Action) (ratelimit.Result, error) {
	level, err := s.GetUserRiskLevel(user)
	if err != nil {
		return ratelimit.Result{Allowed: true}, nil // 降级放行
	}
	return s.limiter.Allow(ctx, fmt.Sprint(user.ID), action, toRatelimitLevel(level))
}

// RecordRiskEvent 记录一次风险事件（举报成立、命中敏感词等）
// ttl: 该事件计入风险分的有效期
func (s *riskService) RecordRiskEvent(userID uint, eventType, detail string, ttl time.Duration) error {
	record := &po.UserRiskRecord{
		UserID:      userID,
		EventType:   eventType,
		EventDetail: detail,
		ExpireAt:    time.Now().Add(ttl),
	}
	return s.repo.AddRiskRecord(record)
}

// WriteAuditLog 写入操作审计日志
func (s *riskService) WriteAuditLog(operatorID uint, action po.AuditActionType,
	targetType string, targetID uint, before, after, reason, ip string) error {
	log := &po.AuditLog{
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
