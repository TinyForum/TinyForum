package request

// CreateReportRequest 创建举报的请求参数
type CreateReportRequest struct {
	TargetID        uint   `json:"target_id" binding:"required"`                                                       // 被举报对象ID
	TargetType      string `json:"target_type" binding:"required,oneof=post comment user"`                             // 被举报对象类型
	Type            string `json:"type" binding:"omitempty,oneof=spam offensive illegal misinformation privacy other"` // 举报类型
	Reason          string `json:"reason" binding:"required,min=2,max=500"`                                            // 举报理由
	ContentSnapshot string `json:"content_snapshot" binding:"omitempty,max=2000"`                                      // 被举报内容快照
	IsAnonymous     bool   `json:"is_anonymous"`                                                                       // 是否匿名举报
	ReporterIP      string `json:"-"`                                                                                  // 举报人IP（仅服务端填充，不接收客户端传入）
}

// HandleReportRequest 处理举报的请求参数
type HandleReportRequest struct {
	Status       string `json:"status" binding:"required,oneof=resolved rejected"` // 处理结果：通过/驳回
	HandleNote   string `json:"handle_note" binding:"required,min=1,max=500"`      // 处理备注
	RejectReason string `json:"reject_reason" binding:"omitempty,max=500"`          // 驳回原因（可选）
}
