package reports

import (
	"context"
	"errors"
	"time"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// toReportVO 将 do.Report 转换为视图对象（脱敏：匿名举报不返回举报人信息）
func toReportVO(r *do.Report) *vo.ReportVO {
	if r == nil {
		return nil
	}
	v := &vo.ReportVO{
		ID:              r.ID,
		TargetID:        r.TargetID,
		TargetType:      string(r.TargetType),
		Type:            string(r.Type),
		Reason:          r.Reason,
		Status:          string(r.Status),
		HandleNote:      r.HandleNote,
		HandleAt:        r.HandleAt,
		ContentSnapshot: r.ContentSnapshot,
		IsAnonymous:     r.IsAnonymous,
		Priority:        r.Priority,
	}
	if !r.IsAnonymous && r.Reporter != nil {
		v.Reporter = &vo.UserPublicVO{
			ID:        r.Reporter.ID,
			Name:      r.Reporter.Username,
			AvatarUrl: r.Reporter.AvatarUrl,
		}
	}
	if r.Handler != nil {
		v.Handler = &vo.UserPublicVO{
			ID:        r.Handler.ID,
			Name:      r.Handler.Username,
			AvatarUrl: r.Handler.AvatarUrl,
		}
	}
	return v
}

// GetByID 获取单条举报详情
func (s *reportsService) GetByID(ctx context.Context, id uint) (*vo.ReportVO, error) {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithDetail("举报不存在")
		}
		return nil, err
	}
	return toReportVO(report), nil
}

// Handle 处理一条举报（通过/驳回），写入处理人、备注与处理时间
func (s *reportsService) Handle(ctx context.Context, id uint, req *request.HandleReportRequest, handlerID uint) (*vo.ReportVO, error) {
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound.WithDetail("举报不存在")
		}
		return nil, err
	}
	if report.Status != do.ReportPending {
		return nil, apperrors.ErrInvalidRequest.WithDetail("该举报已被处理")
	}

	now := time.Now()
	report.Status = do.ReportStatus(req.Status)
	report.HandlerID = &handlerID
	report.HandleNote = req.HandleNote
	if req.RejectReason != "" {
		report.HandleNote = req.HandleNote + "；驳回原因：" + req.RejectReason
	}
	report.HandleAt = &now

	if err := s.reportRepo.Update(ctx, report); err != nil {
		return nil, err
	}
	return toReportVO(report), nil
}
