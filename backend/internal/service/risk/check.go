// internal/service/risk/risk_service.go

package risk

import (
	"context"
	"fmt"
	"log"
	"time"
	"tiny-forum/internal/infra/ratelimit"
	"tiny-forum/internal/model"
)

// type riskService struct {
// 	repo    RiskRepository
// 	limiter ratelimit.Limiter
// }

// // RiskRepository 接口需要扩展
// type RiskRepository interface {
// 	// 现有方法
// 	CountActiveRiskEvents(userID uint) (int, error)
// 	AddRiskRecord(record *model.UserRiskRecord) error
// 	CreateAuditLog(log *model.AuditLog) error

// 	// 新增方法
// 	CountActiveRiskEventsByIP(ip string) (int, error)
// 	AddIPRiskRecord(record *model.IPRiskRecord) error
// 	IsIPBlocked(ip string) (bool, error)
// }

// GetAnonymousRiskLevel 获取匿名用户（未登录）的风险等级
func (s *riskService) GetAnonymousRiskLevel(ip string) (model.RiskLevel, error) {
	// 可选：检查IP是否被封锁
	// isBlocked, err := s.repo.IsIPBlocked(ip)
	// if err != nil {
	//     return model.RiskLevelNormal, err
	// }
	// if isBlocked {
	//     return model.RiskLevelBlocked, nil
	// }

	// 统计IP的活跃风险事件数
	activeEvents, err := s.repo.CountActiveRiskEventsByIP(ip)
	if err != nil {
		return model.RiskLevelNormal, err
	}

	// 活跃风险事件 >= 3 则限制
	if activeEvents >= 3 {
		return model.RiskLevelRestrict, nil
	}

	// 匿名用户默认为正常等级
	return model.RiskLevelNormal, nil
}

// CheckRateLimitByIP 检查匿名用户（未登录）操作频率是否超限
func (s *riskService) CheckRateLimitByIP(ctx context.Context, ip string, action ratelimit.Action) (ratelimit.Result, error) {
	log.Printf("[RateLimit] Anonymous user IP: %s, action: %s", ip, action)

	// 未登录用户使用固定的 IP 限流规则
	quota := getAnonymousQuota(action)

	// 使用 IP 作为标识符，固定使用 normal 等级
	identifier := fmt.Sprintf("ip:%s", ip)

	return s.limiter.Allow(ctx, identifier, action, ratelimit.RiskNormal, quota)
}

func getAnonymousQuota(action ratelimit.Action) ratelimit.Quota {
	switch action {
	case ratelimit.ActionLogin:
		return ratelimit.Quota{Limit: 5, Window: 5 * time.Minute} // 5分钟内最多5次
	case ratelimit.ActionRegister:
		return ratelimit.Quota{Limit: 3, Window: 1 * time.Hour} // 1小时内最多3次
	case ratelimit.ActionGetPost, ratelimit.ActionGetComment:
		return ratelimit.Quota{Limit: 100, Window: 1 * time.Minute} // 读取操作宽松
	default:
		return ratelimit.Quota{Limit: 10, Window: time.Minute} // 默认：每分钟10次
	}
}

// RecordRiskEventByIP 记录基于IP的风险事件
func (s *riskService) RecordRiskEventByIP(ip, eventType, detail string, ttl time.Duration) error {
	record := &model.IPRiskRecord{
		IP:          ip,
		EventType:   eventType,
		EventDetail: detail,
		ExpireAt:    time.Now().Add(ttl),
	}
	return s.repo.AddIPRiskRecord(record)
}

// WriteAuditLogByIP 写入基于IP的操作审计日志
func (s *riskService) WriteAuditLogByIP(ip string, action model.AuditActionType,
	targetType string, targetID uint, before, after, reason string) error {
	log := &model.AuditLog{
		OperatorIP: ip, // 需要在 AuditLog 模型中添加 OperatorIP 字段
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
