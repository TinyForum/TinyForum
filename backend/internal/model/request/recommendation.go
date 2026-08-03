package request

// RecordBehaviorRequest 记录用户行为请求
type RecordBehaviorRequest struct {
	TargetID    uint   `json:"target_id" binding:"required"`
	TargetType  string `json:"target_type" binding:"required"`
	BehaviorType string `json:"behavior_type" binding:"required"`
	Value       float64 `json:"value"`
	SessionID   string  `json:"session_id"`
	ContextJSON string  `json:"context_json"`
}

// RecommendationQuery 推荐查询参数
type RecommendationQuery struct {
	Page     int    `form:"page" json:"page"`
	PageSize int    `form:"page_size" json:"page_size"`
	Strategy string `form:"strategy" json:"strategy"`
	Cursor   string `form:"cursor" json:"cursor"`
	BoardID  uint   `form:"board_id" json:"board_id"`
}

// RecommendationFeedbackRequest 推荐反馈请求
type RecommendationFeedbackRequest struct {
	CreationID   uint   `json:"creation_id" binding:"required"`
	FeedbackType string `json:"feedback_type" binding:"required"`
	SourceType   string `json:"source_type"`
	Position     int    `json:"position"`
	SessionID    string `json:"session_id"`
}

// BatchFeedbackRequest 批量反馈请求
type BatchFeedbackRequest struct {
	Feedbacks []RecommendationFeedbackRequest `json:"feedbacks" binding:"required"`
}
