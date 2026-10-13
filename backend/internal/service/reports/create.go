package reports

import (
	"context"
	"fmt"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	apperrors "tiny-forum/pkg/errors"
)

// CreateReport 创建举报：对同一内容的重复举报会记录在案，由管理员人工处理并实现状态流转。
func (s *reportsService) CreateReport(ctx context.Context, reporterID uint, req *request.CreateReportRequest) (*vo.ReportVO, error) {
	if reporterID == 0 {
		return nil, apperrors.ErrUnauthorized
	}

	reportType := do.ReportType(req.Type)
	if reportType == "" {
		reportType = do.ReportTypeOther
	}

	report := &do.Report{
		ReporterID:      reporterID,
		TargetID:        req.TargetID,
		TargetType:      do.ReportTargetType(req.TargetType),
		Type:            reportType,
		Reason:          req.Reason,
		Status:          do.ReportPending,
		ContentSnapshot: req.ContentSnapshot,
		ReporterIP:      req.ReporterIP,
		IsAnonymous:     req.IsAnonymous,
		// HandlerID / HandleNote / HandleAt 由管理员处理时填充
	}

	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("创建举报失败: %w", err)
	}

	// 举报聚合：当同一目标待处理举报达到阈值时，创建审核任务
	if s.checkSvc != nil {
		_, _ = s.checkSvc.HandleReportAggregate(do.AuditTargetType(report.TargetType), report.TargetID)
	}

	// 返回 VO 而非 do，避免泄漏 ReporterIP/HandlerID 等内部字段
	return &vo.ReportVO{
		ID:              report.ID,
		TargetID:        report.TargetID,
		TargetType:      string(report.TargetType),
		Type:            string(report.Type),
		Reason:          report.Reason,
		Status:          string(report.Status),
		ContentSnapshot: report.ContentSnapshot,
		IsAnonymous:     report.IsAnonymous,
		Priority:        report.Priority,
	}, nil
}
