package user

import (
	"context"

	"tiny-forum/pkg/logger"

	"go.uber.org/zap"
)

// logAudit 记录管理员敏感操作的审计日志（结构化写入 Zap）
func (s *userService) logAudit(_ context.Context, operatorID, targetID uint, action, detail string) error {
	logger.Info("管理员操作审计",
		zap.Uint("operator_id", operatorID),
		zap.Uint("target_id", targetID),
		zap.String("action", action),
		zap.String("detail", detail))
	return nil
}
