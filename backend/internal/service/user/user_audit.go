package user

import "context"

// logAudit 记录审计日志（可扩展）
func (s *UserService) logAudit(ctx context.Context, operatorID, targetID uint, action, detail string) {
	// 如果项目有日志库或审计表，在此实现
	// 例如：s.logger.Info("audit_log", "operator", operatorID, "target", targetID, "action", action, "detail", detail)
}
