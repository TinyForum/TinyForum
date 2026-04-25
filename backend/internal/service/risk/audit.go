package risk

import "tiny-forum/internal/model"

// GetAuditLogs 查询审计日志（供 handler 调用）
func (s *riskService) GetAuditLogs(targetType string, targetID uint, limit int) ([]model.AuditLog, error) {
	return s.repo.ListAuditLogs(targetType, targetID, limit)
}
