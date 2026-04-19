package risk

import (
	"context"
	"time"
	"tiny-forum/internal/model"
	riskrepo "tiny-forum/internal/repository/risk"
	"tiny-forum/pkg/ratelimit"
)

// RiskService 风控核心服务
type RiskService struct {
	repo    *riskrepo.RiskRepository
	limiter *ratelimit.Limiter
}

func NewRiskService(repo *riskrepo.RiskRepository, limiter *ratelimit.Limiter) *RiskService {
	return &RiskService{repo: repo, limiter: limiter}
}

// GetUserRiskLevel 计算用户当前风险等级
// 规则（优先级从高到低）：
//  1. IsBlocked → blocked（调用方直接拦截，无需走此函数）
//  2. 活跃风险事件 >= 3 → restrict
//  3. score < 50 或 注册不足7天 → observe
//  4. 其他 → normal
func (s *RiskService) GetUserRiskLevel(user *model.User) (model.RiskLevel, error) {
	if user.IsBlocked {
		return model.RiskLevelBlocked, nil
	}

	activeEvents, err := s.repo.CountActiveRiskEvents(user.ID)
	if err != nil {
		return model.RiskLevelNormal, err
	}
	if activeEvents >= 3 {
		return model.RiskLevelRestrict, nil
	}

	isNewUser := time.Since(user.CreatedAt) < 7*24*time.Hour
	isLowScore := user.Score < 50
	if isNewUser || isLowScore {
		return model.RiskLevelObserve, nil
	}

	return model.RiskLevelNormal, nil
}

// toRatelimitLevel 将 model.RiskLevel 转换为 ratelimit 包的类型（避免循环依赖）
func toRatelimitLevel(level model.RiskLevel) ratelimit.RiskLevel {
	switch level {
	case model.RiskLevelRestrict:
		return ratelimit.RiskRestrict
	case model.RiskLevelObserve:
		return ratelimit.RiskObserve
	default:
		return ratelimit.RiskNormal
	}
}

// CheckRateLimit 检查用户操作频率是否超限
func (s *RiskService) CheckRateLimit(ctx context.Context, user *model.User, action ratelimit.Action) (ratelimit.Result, error) {
	level, err := s.GetUserRiskLevel(user)
	if err != nil {
		return ratelimit.Result{Allowed: true}, nil // 降级放行
	}
	return s.limiter.Allow(ctx, user.ID, action, toRatelimitLevel(level))
}

// RecordRiskEvent 记录一次风险事件（举报成立、命中敏感词等）
// ttl: 该事件计入风险分的有效期
func (s *RiskService) RecordRiskEvent(userID uint, eventType, detail string, ttl time.Duration) error {
	record := &model.UserRiskRecord{
		UserID:      userID,
		EventType:   eventType,
		EventDetail: detail,
		ExpireAt:    time.Now().Add(ttl),
	}
	return s.repo.AddRiskRecord(record)
}

// WriteAuditLog 写入操作审计日志
func (s *RiskService) WriteAuditLog(operatorID uint, action model.AuditActionType,
	targetType string, targetID uint, before, after, reason, ip string) error {
	log := &model.AuditLog{
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
