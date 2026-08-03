package vo

import "time"

// RecommendationItem 推荐内容项
type RecommendationItem struct {
	CreationID   uint    `json:"creation_id"`
	Title        string  `json:"title"`
	Summary      string  `json:"summary"`
	CoverURL     string  `json:"cover_url"`
	AuthorID     uint    `json:"author_id"`
	AuthorName   string  `json:"author_name"`
	AuthorAvatar string  `json:"author_avatar"`
	ViewCount    int     `json:"view_count"`
	LikeCount    int     `json:"like_count"`
	CommentCount int     `json:"comment_count"`
	BoardID      uint    `json:"board_id"`
	BoardName    string  `json:"board_name"`
	Score        float64 `json:"score"`
	Reason       string  `json:"reason"`
	CreatedAt    string  `json:"created_at"`
}

// RecommendationResponse 推荐响应
type RecommendationResponse struct {
	Items      []RecommendationItem `json:"items"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	HasMore    bool                 `json:"has_more"`
	SessionID  string               `json:"session_id"`
	Strategy   string               `json:"strategy"`
	GeneratedAt time.Time           `json:"generated_at"`
}

// BehaviorEventVO 行为事件视图
type BehaviorEventVO struct {
	ID           uint   `json:"id"`
	TargetID     uint   `json:"target_id"`
	TargetType   string `json:"target_type"`
	BehaviorType string `json:"behavior_type"`
	CreatedAt    string `json:"created_at"`
}

// FeedbackResult 反馈结果
type FeedbackResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// UserInterestProfileVO 用户兴趣画像视图
type UserInterestProfileVO struct {
	ActiveTags  []string           `json:"active_tags"`
	TagWeights  map[string]float64 `json:"tag_weights"`
	UpdatedAt   string             `json:"updated_at"`
}

// --- 管理员分析 VO ---

// RecOverviewStats 推荐系统概览统计
type RecOverviewStats struct {
	TotalBehaviors    int64   `json:"total_behaviors"`
	TotalFeedbacks    int64   `json:"total_feedbacks"`
	TodayBehaviors    int64   `json:"today_behaviors"`
	TodayImpressions  int64   `json:"today_impressions"`
	ClickRate         float64 `json:"click_rate"`
	DismissRate       float64 `json:"dismiss_rate"`
	UserCount         int64   `json:"user_count"`
	ContentCount      int64   `json:"content_count"`
	AvgQualityScore   float64 `json:"avg_quality_score"`
}

// BehaviorDistribution 行为分布
type BehaviorDistribution struct {
	BehaviorType string  `json:"behavior_type"`
	Count        int64   `json:"count"`
	Ratio        float64 `json:"ratio"`
}

// BehaviorStats 行为统计
type BehaviorStats struct {
	Distribution []BehaviorDistribution `json:"distribution"`
	DailyTrend   []TrendPoint           `json:"daily_trend"`
	TopTargets   []TopTargetItem        `json:"top_targets"`
}

// TrendPoint 趋势数据点
type TrendPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// TopTargetItem 热门目标内容
type TopTargetItem struct {
	TargetID     uint   `json:"target_id"`
	Title        string `json:"title"`
	BehaviorCount int64  `json:"behavior_count"`
}

// UserAnalysis 用户分析数据
type UserAnalysis struct {
	TotalTrackedUsers  int64             `json:"total_tracked_users"`
	ActiveUsersToday   int64             `json:"active_users_today"`
	AvgBehaviorsPerUser float64          `json:"avg_behaviors_per_user"`
	TopActiveUsers     []ActiveUserItem  `json:"top_active_users"`
	TopInterestTags    []TagWeightItem   `json:"top_interest_tags"`
}

// ActiveUserItem 活跃用户项
type ActiveUserItem struct {
	UserID         uint   `json:"user_id"`
	Username       string `json:"username"`
	Nickname       string `json:"nickname"`
	BehaviorCount  int64  `json:"behavior_count"`
	TopBehavior    string `json:"top_behavior"`
}

// TagWeightItem 标签权重项
type TagWeightItem struct {
	TagID    uint    `json:"tag_id"`
	TagName  string  `json:"tag_name"`
	Weight   float64 `json:"weight"`
	UserCount int64  `json:"user_count"`
}

// ContentPerformance 内容表现分析
type ContentPerformance struct {
	TopHotContent     []ContentPerfItem `json:"top_hot_content"`
	TopQualityContent []ContentPerfItem `json:"top_quality_content"`
	ContentCountByBoard []BoardContentCount `json:"content_count_by_board"`
	QualityDistribution []ScoreBucket   `json:"quality_distribution"`
}

// ContentPerfItem 内容表现项
type ContentPerfItem struct {
	CreationID    uint    `json:"creation_id"`
	Title         string  `json:"title"`
	ViewCount     int     `json:"view_count"`
	LikeCount     int     `json:"like_count"`
	HotScore      float64 `json:"hot_score"`
	QualityScore  float64 `json:"quality_score"`
}

// BoardContentCount 板块内容计数
type BoardContentCount struct {
	BoardID uint   `json:"board_id"`
	BoardName string `json:"board_name"`
	Count    int64  `json:"count"`
}

// ScoreBucket 分数区间
type ScoreBucket struct {
	Range string `json:"range"`
	Count int64  `json:"count"`
}

// --- 风控关联分析 VO ---

// RiskAnalysis 风控关联分析
type RiskAnalysis struct {
	TotalRiskUsers      int64          `json:"total_risk_users"`
	DangerLevelUsers    int64          `json:"danger_level_users"`
	TotalViolations     int64          `json:"total_violations"`
	PendingViolations   int64          `json:"pending_violations"`
	TotalBans           int64          `json:"total_bans"`
	ViolationDistribution []ViolationTypeItem `json:"violation_distribution"`
	RiskUserBehaviors   []RiskUserBehavior    `json:"risk_user_behaviors"`
	TopReportedContent  []ReportedContentItem `json:"top_reported_content"`
}

// ViolationTypeItem 违规类型分布
type ViolationTypeItem struct {
	ViolationType string  `json:"violation_type"`
	Count         int64   `json:"count"`
	Ratio         float64 `json:"ratio"`
}

// RiskUserBehavior 风控用户行为关联
type RiskUserBehavior struct {
	UserID        uint   `json:"user_id"`
	Username      string `json:"username"`
	RiskLevel     string `json:"risk_level"`
	ViolationCount int64 `json:"violation_count"`
	BehaviorCount int64  `json:"behavior_count"`
	TopBehavior   string `json:"top_behavior"`
}

// ReportedContentItem 被举报内容
type ReportedContentItem struct {
	CreationID   uint   `json:"creation_id"`
	Title        string `json:"title"`
	ReportCount  int64  `json:"report_count"`
	RiskLevel    string `json:"risk_level"`
}

// --- 用户综合分析 VO（管理后台-用户分析标签） ---

// ComprehensiveUserAnalysis 用户综合分析
type ComprehensiveUserAnalysis struct {
	Overview       UserAnalysisOverview    `json:"overview"`
	TagDistribution []TagUserDistribution   `json:"tag_distribution"`
	UserBehaviorPatterns []UserBehaviorPattern `json:"user_behavior_patterns"`
	SimilarUserGroups []SimilarUserGroup    `json:"similar_user_groups"`
	RiskUserList      []UserRiskProfile     `json:"risk_user_list"`
}

// UserAnalysisOverview 用户分析概览
type UserAnalysisOverview struct {
	TotalUsers            int64   `json:"total_users"`
	TotalBehaviorRecords  int64   `json:"total_behavior_records"`
	UsersWithProfile      int64   `json:"users_with_profile"`
	AvgTagsPerUser        float64 `json:"avg_tags_per_user"`
	AvgBehaviorsPerUser24h float64 `json:"avg_behaviors_per_user_24h"`
	SharedUsersCount      int64   `json:"shared_users_count"`      // 有行为记录的用户数
	RiskUserCount         int64   `json:"risk_user_count"`
	ViolationUserCount    int64   `json:"violation_user_count"`
}

// TagUserDistribution 标签用户分布
type TagUserDistribution struct {
	TagID         uint    `json:"tag_id"`
	TagName       string  `json:"tag_name"`
	UserCount     int64   `json:"user_count"`
	AvgWeight     float64 `json:"avg_weight"`
	PostCount     int64   `json:"post_count"`
}

// UserBehaviorPattern 用户行为模式（单个用户维度）
type UserBehaviorPattern struct {
	UserID        uint                    `json:"user_id"`
	Username      string                  `json:"username"`
	Nickname      string                  `json:"nickname"`
	Avatar        string                  `json:"avatar"`
	TotalBehaviors int64                  `json:"total_behaviors"`
	BehaviorBreakdown map[string]int64    `json:"behavior_breakdown"`
	ActiveTags    []string                `json:"active_tags"`
	LastActiveAt  string                  `json:"last_active_at"`
	RiskLevel     string                  `json:"risk_level"`
	ViolationCount int64                 `json:"violation_count"`
}

// SimilarUserGroup 相似用户群组
type SimilarUserGroup struct {
	SeedUserID   uint              `json:"seed_user_id"`
	SeedUsername string            `json:"seed_username"`
	SimilarUsers []SimilarUserItem `json:"similar_users"`
}

// SimilarUserItem 相似用户项
type SimilarUserItem struct {
	UserID        uint    `json:"user_id"`
	Username      string  `json:"username"`
	Nickname      string  `json:"nickname"`
	Avatar        string  `json:"avatar"`
	SimilarityScore float64 `json:"similarity_score"`
	SharedTags    []string `json:"shared_tags"`
	CommonBehavior string  `json:"common_behavior"`
}

// UserRiskProfile 用户风险画像
type UserRiskProfile struct {
	UserID         uint   `json:"user_id"`
	Username       string `json:"username"`
	RiskLevel      string `json:"risk_level"`
	ViolationCount int64  `json:"violation_count"`
	BehaviorCount  int64  `json:"behavior_count"`
	LastViolationType string `json:"last_violation_type"`
	LastViolationAt   string `json:"last_violation_at"`
	IsBanned       bool   `json:"is_banned"`
}

