package vo

// ========================================================================
// 管理后台 — 用户分析（6 大模块）VO
// ========================================================================

// ── 1. 用户概览 (Overview) ──

type UAOverview struct {
	TotalUsers          int64   `json:"total_users"`
	DAU                 int64   `json:"dau"`
	WAU                 int64   `json:"wau"`
	MAU                 int64   `json:"mau"`
	Stickiness          float64 `json:"stickiness"`
	NewUsersToday       int64   `json:"new_users_today"`
	NewUsersWeek        int64   `json:"new_users_week"`
	RetentionDay7       float64 `json:"retention_day7"`
	RetentionDay30      float64 `json:"retention_day30"`
	AvgDailyActions     float64 `json:"avg_daily_actions"`
	AvgSessionDuration  float64 `json:"avg_session_duration"`
	TrendPoints         []UATrendPoint `json:"trend_points"`
}

type UATrendPoint struct {
	Date    string `json:"date"`
	DAU     int64  `json:"dau"`
	NewUser int64  `json:"new_user"`
	Actions int64  `json:"actions"`
}

// ── 2. 用户画像 (Profiles) ──

type UAProfile struct {
	UserID            uint              `json:"user_id"`
	Username          string            `json:"username"`
	Nickname          string            `json:"nickname"`
	Avatar            string            `json:"avatar"`
	Email             string            `json:"email"`
	JoinedAt          string            `json:"joined_at"`
	LastActiveAt      string            `json:"last_active_at"`
	EngagementTier    string            `json:"engagement_tier"`
	TotalActions      int64             `json:"total_actions"`
	ActiveDays        int64             `json:"active_days"`
	ActiveTags        []string          `json:"active_tags"`
	TopBehaviors      map[string]int64  `json:"top_behaviors"`
	ContentCount      int64             `json:"content_count"`
	CommentCount      int64             `json:"comment_count"`
	LikeReceived      int64             `json:"like_received"`
	RiskLevel         string            `json:"risk_level"`
	ViolationCount    int64             `json:"violation_count"`
	IsBanned          bool              `json:"is_banned"`
}

type UAProfileList struct {
	Profiles []UAProfile `json:"profiles"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// ── 3. 用户分群 (Segments) ──

type UASegment struct {
	SegmentName     string  `json:"segment_name"`
	UserCount       int64   `json:"user_count"`
	Percentage      float64 `json:"percentage"`
	AvgActions      float64 `json:"avg_actions"`
	AvgActiveDays   float64 `json:"avg_active_days"`
	RetentionRate   float64 `json:"retention_rate"`
	TopTags         []string `json:"top_tags"`
	RiskRatio        float64 `json:"risk_ratio"`
}

type UASegments struct {
	Segments     []UASegment `json:"segments"`
	TotalUsers   int64       `json:"total_users"`
}

// ── 4. 行为分析 (Behavior) ──

type UABehaviorFunnel struct {
	Steps []UAFunnelStep `json:"steps"`
}

type UAFunnelStep struct {
	StepName      string  `json:"step_name"`
	UserCount     int64   `json:"user_count"`
	Conversion    float64 `json:"conversion"`
	DropOff       float64 `json:"drop_off"`
}

type UABehavior struct {
	Funnel          UABehaviorFunnel       `json:"funnel"`
	EventDistribution []UAEventItem        `json:"event_distribution"`
	HourlyHeatmap   []UAHourlyItem         `json:"hourly_heatmap"`
	TopEventUsers   []UAEventUserItem      `json:"top_event_users"`
}

type UAEventItem struct {
	EventType string  `json:"event_type"`
	Count     int64   `json:"count"`
	UniqueUsers int64 `json:"unique_users"`
	Ratio     float64 `json:"ratio"`
}

type UAHourlyItem struct {
	Hour  int   `json:"hour"`
	Count int64 `json:"count"`
}

type UAEventUserItem struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Count    int64  `json:"count"`
	EventType string `json:"event_type"`
}

// ── 5. 同期群分析 (Cohorts) ──

type UACohort struct {
	CohortLabel    string  `json:"cohort_label"`
	InitialUsers   int64   `json:"initial_users"`
	RetentionW1    float64 `json:"retention_w1"`
	RetentionW2    float64 `json:"retention_w2"`
	RetentionW3    float64 `json:"retention_w3"`
	RetentionW4    float64 `json:"retention_w4"`
	RetentionW8    float64 `json:"retention_w8"`
}

type UACohorts struct {
	Cohorts    []UACohort `json:"cohorts"`
	WeekLabels []string   `json:"week_labels"`
	MaxWeeks   int        `json:"max_weeks"`
}

// ── 6. 风险评估 (Risk) ──

type UARiskScoreDist struct {
	Range string `json:"range"`
	Count int64  `json:"count"`
}

type UARisk struct {
	TotalRiskUsers      int64             `json:"total_risk_users"`
	HighRiskCount       int64             `json:"high_risk_count"`
	NewRisksToday       int64             `json:"new_risks_today"`
	PendingReviews      int64             `json:"pending_reviews"`
	ScoreDistribution   []UARiskScoreDist `json:"score_distribution"`
	ViolationTrend      []UATrendPoint    `json:"violation_trend"`
	TopRiskyUsers       []UARiskyUserItem `json:"top_risky_users"`
	FlaggedContentQueue []UAFlaggedItem   `json:"flagged_content_queue"`
}

type UARiskyUserItem struct {
	UserID         uint   `json:"user_id"`
	Username       string `json:"username"`
	RiskLevel      string `json:"risk_level"`
	ViolationCount int64  `json:"violation_count"`
	BehaviorCount  int64  `json:"behavior_count"`
	LastFlagReason string `json:"last_flag_reason"`
	IsBanned       bool   `json:"is_banned"`
}

type UAFlaggedItem struct {
	CreationID   uint   `json:"creation_id"`
	Title        string `json:"title"`
	AuthorName   string `json:"author_name"`
	ReportCount  int64  `json:"report_count"`
	FlagReason   string `json:"flag_reason"`
	CreatedAt    string `json:"created_at"`
}
