package bo

import (
	"time"
	"tiny-forum/internal/model/do"
)

type ListReportBO struct {
	ID         uint
	ReporterID uint
	TargetID   uint
	TargetType string
	Type       do.ReportType
	Reason     string
	Status     do.ReportStatus
	HandlerID  *uint
	HandleNote string     //
	HandleAt   *time.Time // 处理时间

	// 以下为扩展推荐字段（可根据需要启用）
	ContentSnapshot string // 内容片段
	ReporterIP      string // 举报IP
	IsAnonymous     bool   //
	Priority        int8
}
