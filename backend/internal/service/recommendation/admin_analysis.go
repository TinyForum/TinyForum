package recommendation

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"time"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
	"tiny-forum/pkg/logger"

	"go.uber.org/zap"
)

// GetOverviewStats 获取推荐系统概览统计
func (s *recommendationService) GetOverviewStats(ctx context.Context) (*vo.RecOverviewStats, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	var totalBehaviors int64
	s.db.WithContext(ctx).Model(&do.UserBehaviorEvent{}).Count(&totalBehaviors)

	var totalFeedbacks int64
	s.db.WithContext(ctx).Model(&do.RecommendationFeedback{}).Count(&totalFeedbacks)

	var todayBehaviors int64
	s.db.WithContext(ctx).Model(&do.UserBehaviorEvent{}).
		Where("created_ts >= ?", todayStart.UnixMilli()).Count(&todayBehaviors)

	var todayImpressions int64
	s.db.WithContext(ctx).Model(&do.RecommendationFeedback{}).
		Where("feedback_type = ? AND created_at >= ?", do.FeedbackImpression, todayStart).Count(&todayImpressions)

	var clickCount int64
	s.db.WithContext(ctx).Model(&do.RecommendationFeedback{}).
		Where("feedback_type = ?", do.FeedbackClick).Count(&clickCount)

	var dismissCount int64
	s.db.WithContext(ctx).Model(&do.RecommendationFeedback{}).
		Where("feedback_type IN ?", []do.FeedbackType{do.FeedbackDismiss, do.FeedbackNotInterested}).Count(&dismissCount)

	clickRate := 0.0
	if totalFeedbacks > 0 {
		clickRate = math.Round(float64(clickCount)/float64(totalFeedbacks)*10000) / 100
	}

	dismissRate := 0.0
	if totalFeedbacks > 0 {
		dismissRate = math.Round(float64(dismissCount)/float64(totalFeedbacks)*10000) / 100
	}

	var userCount int64
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE deleted_at IS NULL",
	).Scan(&userCount)

	var contentCount int64
	s.db.WithContext(ctx).Model(&do.ContentFeature{}).Count(&contentCount)

	var avgQuality float64
	s.db.WithContext(ctx).Raw(
		"SELECT COALESCE(AVG(quality_score), 0) FROM content_features WHERE deleted_at IS NULL",
	).Scan(&avgQuality)

	return &vo.RecOverviewStats{
		TotalBehaviors:   totalBehaviors,
		TotalFeedbacks:   totalFeedbacks,
		TodayBehaviors:   todayBehaviors,
		TodayImpressions: todayImpressions,
		ClickRate:        clickRate,
		DismissRate:      dismissRate,
		UserCount:        userCount,
		ContentCount:     contentCount,
		AvgQualityScore:  avgQuality,
	}, nil
}

// GetBehaviorStats 获取行为统计（分布 + 趋势 + 热门目标）
func (s *recommendationService) GetBehaviorStats(ctx context.Context, days int) (*vo.BehaviorStats, error) {
	if days < 1 {
		days = 7
	}
	since := time.Now().AddDate(0, 0, -days)

	result := &vo.BehaviorStats{}

	// 行为分布
	type behaviorRow struct {
		BehaviorType string
		Count        int64
	}
	var rows []behaviorRow
	err := s.db.WithContext(ctx).Raw(`
		SELECT behavior_type, COUNT(*) as count
		FROM user_behavior_events
		WHERE deleted_at IS NULL AND created_at >= ?
		GROUP BY behavior_type
		ORDER BY count DESC
	`, since).Scan(&rows).Error
	if err != nil {
		logger.Error("推荐系统：查询行为分布失败", zap.Error(err))
		return nil, err
	}

	var total int64
	for _, r := range rows {
		total += r.Count
	}
	for _, r := range rows {
		ratio := 0.0
		if total > 0 {
			ratio = math.Round(float64(r.Count)/float64(total)*10000) / 100
		}
		result.Distribution = append(result.Distribution, vo.BehaviorDistribution{
			BehaviorType: r.BehaviorType,
			Count:        r.Count,
			Ratio:        ratio,
		})
	}

	// 每日趋势
	type trendRow struct {
		Date  string
		Count int64
	}
	var trends []trendRow
	s.db.WithContext(ctx).Raw(`
		SELECT TO_CHAR(DATE(created_at), 'YYYY-MM-DD') as date, COUNT(*) as count
		FROM user_behavior_events
		WHERE deleted_at IS NULL AND created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY date
	`, since).Scan(&trends)

	for _, t := range trends {
		result.DailyTrend = append(result.DailyTrend, vo.TrendPoint{
			Date:  t.Date,
			Count: t.Count,
		})
	}

	// 热门目标内容
	type topRow struct {
		TargetID      uint
		BehaviorCount int64
	}
	var tops []topRow
	s.db.WithContext(ctx).Raw(`
		SELECT target_id, COUNT(*) as behavior_count
		FROM user_behavior_events
		WHERE deleted_at IS NULL AND created_at >= ?
		GROUP BY target_id
		ORDER BY behavior_count DESC
		LIMIT 10
	`, since).Scan(&tops)

	targetIDs := make([]uint, len(tops))
	for i, t := range tops {
		targetIDs[i] = t.TargetID
	}

	type titleRow struct {
		ID    uint
		Title string
	}
	var titles []titleRow
	if len(targetIDs) > 0 {
		s.db.WithContext(ctx).Raw(`
			SELECT id, title FROM creations WHERE id IN ? AND deleted_at IS NULL
		`, targetIDs).Scan(&titles)
	}
	titleMap := make(map[uint]string, len(titles))
	for _, t := range titles {
		titleMap[t.ID] = t.Title
	}

	for _, t := range tops {
		result.TopTargets = append(result.TopTargets, vo.TopTargetItem{
			TargetID:      t.TargetID,
			Title:         titleMap[t.TargetID],
			BehaviorCount: t.BehaviorCount,
		})
	}

	return result, nil
}

// GetUserAnalysis 获取用户分析数据
func (s *recommendationService) GetUserAnalysis(ctx context.Context, days int) (*vo.UserAnalysis, error) {
	if days < 1 {
		days = 7
	}
	since := time.Now().AddDate(0, 0, -days)
	todayStart := time.Now().Truncate(24 * time.Hour)

	result := &vo.UserAnalysis{}

	// 被追踪的用户总数
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE deleted_at IS NULL",
	).Scan(&result.TotalTrackedUsers)

	// 今日活跃用户
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE deleted_at IS NULL AND created_at >= ?",
		todayStart,
	).Scan(&result.ActiveUsersToday)

	// 人均行为数
	var totalUsers int64
	var totalBehaviors int64
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE deleted_at IS NULL AND created_at >= ?",
		since,
	).Scan(&totalUsers)
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM user_behavior_events WHERE deleted_at IS NULL AND created_at >= ?",
		since,
	).Scan(&totalBehaviors)
	if totalUsers > 0 {
		result.AvgBehaviorsPerUser = math.Round(float64(totalBehaviors)/float64(totalUsers)*10) / 10
	}

	// 最活跃用户
	type activeUserRow struct {
		UserID uint
		Cnt    int64
	}
	var activeUsers []activeUserRow
	s.db.WithContext(ctx).Raw(`
		SELECT user_id, COUNT(*) as cnt
		FROM user_behavior_events
		WHERE deleted_at IS NULL AND created_at >= ?
		GROUP BY user_id
		ORDER BY cnt DESC
		LIMIT 10
	`, since).Scan(&activeUsers)

	userIDs := make([]uint, len(activeUsers))
	for i, u := range activeUsers {
		userIDs[i] = u.UserID
	}

	type userInfoRow struct {
		ID       uint
		Username string
		Nickname string
	}
	var userInfos []userInfoRow
	if len(userIDs) > 0 {
		s.db.WithContext(ctx).Raw(
			"SELECT id, username, nickname FROM users WHERE id IN ? AND deleted_at IS NULL",
			userIDs,
		).Scan(&userInfos)
	}
	userInfoMap := make(map[uint]userInfoRow, len(userInfos))
	for _, u := range userInfos {
		userInfoMap[u.ID] = u
	}

	// 查询每个用户最频繁的行为类型
	type topBehaviorRow struct {
		UserID       uint
		BehaviorType string
	}
	var topBehaviors []topBehaviorRow
	for _, uid := range userIDs {
		var bType string
		s.db.WithContext(ctx).Raw(`
			SELECT behavior_type FROM user_behavior_events
			WHERE user_id = ? AND deleted_at IS NULL AND created_at >= ?
			GROUP BY behavior_type ORDER BY COUNT(*) DESC LIMIT 1
		`, uid, since).Scan(&bType)
		topBehaviors = append(topBehaviors, topBehaviorRow{UserID: uid, BehaviorType: bType})
	}
	topBehaviorMap := make(map[uint]string, len(topBehaviors))
	for _, tb := range topBehaviors {
		topBehaviorMap[tb.UserID] = tb.BehaviorType
	}

	for _, u := range activeUsers {
		info := userInfoMap[u.UserID]
		name := info.Nickname
		if name == "" {
			name = info.Username
		}
		result.TopActiveUsers = append(result.TopActiveUsers, vo.ActiveUserItem{
			UserID:        u.UserID,
			Username:      info.Username,
			Nickname:      name,
			BehaviorCount: u.Cnt,
			TopBehavior:   topBehaviorMap[u.UserID],
		})
	}

	// 最热兴趣标签（聚合所有用户的interest profile）
	type tagWeightAgg struct {
		TagID     uint
		UserCount int64
		Weight    float64
	}
	var tagWeights []tagWeightAgg
	_ = tagWeights // 通过profiles聚合

	// 从user_interest_profiles查找
	type profileRow struct {
		TagWeightsJSON string
	}
	var profiles []profileRow
	s.db.WithContext(ctx).Raw(`
		SELECT tag_weights_json FROM user_interest_profiles WHERE deleted_at IS NULL
	`).Scan(&profiles)

	tagWeightMap := make(map[uint]float64)
	tagUserCountMap := make(map[uint]int64)
	for _, p := range profiles {
		var weights map[string]float64
		if err := json.Unmarshal([]byte(p.TagWeightsJSON), &weights); err != nil {
			continue
		}
		for tidStr, w := range weights {
			tid, err := strconv.Atoi(tidStr)
			if err != nil {
				continue
			}
			tidUint := uint(tid)
			tagWeightMap[tidUint] += w
			tagUserCountMap[tidUint]++
		}
	}

	type tagEntry struct {
		id    uint
		w     float64
		count int64
	}
	var entries []tagEntry
	for tid, w := range tagWeightMap {
		entries = append(entries, tagEntry{tid, w, tagUserCountMap[tid]})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].w > entries[j].w })

	topNTags := 10
	if len(entries) < topNTags {
		topNTags = len(entries)
	}

	tagIDs := make([]uint, topNTags)
	for i := 0; i < topNTags; i++ {
		tagIDs[i] = entries[i].id
	}

	type tagNameRow struct {
		ID   uint
		Name string
	}
	var tagNames []tagNameRow
	if len(tagIDs) > 0 {
		s.db.WithContext(ctx).Raw(
			"SELECT id, name FROM tags WHERE id IN ? AND deleted_at IS NULL", tagIDs,
		).Scan(&tagNames)
	}
	tagNameMap := make(map[uint]string, len(tagNames))
	for _, tn := range tagNames {
		tagNameMap[tn.ID] = tn.Name
	}

	for i := 0; i < topNTags; i++ {
		e := entries[i]
		result.TopInterestTags = append(result.TopInterestTags, vo.TagWeightItem{
			TagID:    e.id,
			TagName:  tagNameMap[e.id],
			Weight:   e.w,
			UserCount: e.count,
		})
	}

	return result, nil
}

// GetContentPerformance 获取内容表现分析
func (s *recommendationService) GetContentPerformance(ctx context.Context) (*vo.ContentPerformance, error) {
	result := &vo.ContentPerformance{}

	// 热度最高内容
	type perfRow struct {
		CreationID   uint
		Title        string
		ViewCount    int
		LikeCount    int
		HotScore     float64
		QualityScore float64
	}
	var hotContents []perfRow
	s.db.WithContext(ctx).Raw(`
		SELECT cf.creation_id, COALESCE(c.title, '') as title,
		       COALESCE(c.view_count, 0) as view_count,
		       COALESCE(c.like_count, 0) as like_count,
		       cf.hot_score, cf.quality_score
		FROM content_features cf
		LEFT JOIN creations c ON c.id = cf.creation_id AND c.deleted_at IS NULL
		WHERE cf.deleted_at IS NULL
		ORDER BY cf.hot_score DESC
		LIMIT 10
	`).Scan(&hotContents)

	for _, hc := range hotContents {
		result.TopHotContent = append(result.TopHotContent, vo.ContentPerfItem{
			CreationID:   hc.CreationID,
			Title:        hc.Title,
			ViewCount:    hc.ViewCount,
			LikeCount:    hc.LikeCount,
			HotScore:     hc.HotScore,
			QualityScore: hc.QualityScore,
		})
	}

	// 质量最高内容
	var qualityContents []perfRow
	s.db.WithContext(ctx).Raw(`
		SELECT cf.creation_id, COALESCE(c.title, '') as title,
		       COALESCE(c.view_count, 0) as view_count,
		       COALESCE(c.like_count, 0) as like_count,
		       cf.hot_score, cf.quality_score
		FROM content_features cf
		LEFT JOIN creations c ON c.id = cf.creation_id AND c.deleted_at IS NULL
		WHERE cf.deleted_at IS NULL
		ORDER BY cf.quality_score DESC
		LIMIT 10
	`).Scan(&qualityContents)

	for _, qc := range qualityContents {
		result.TopQualityContent = append(result.TopQualityContent, vo.ContentPerfItem{
			CreationID:   qc.CreationID,
			Title:        qc.Title,
			ViewCount:    qc.ViewCount,
			LikeCount:    qc.LikeCount,
			HotScore:     qc.HotScore,
			QualityScore: qc.QualityScore,
		})
	}

	// 板块内容分布
	type boardRow struct {
		BoardID   uint
		BoardName string
		Count     int64
	}
	var boardRows []boardRow
	s.db.WithContext(ctx).Raw(`
		SELECT cf.board_id, COALESCE(b.name, '未分类') as board_name, COUNT(*) as count
		FROM content_features cf
		LEFT JOIN boards b ON b.id = cf.board_id AND b.deleted_at IS NULL
		WHERE cf.deleted_at IS NULL
		GROUP BY cf.board_id, b.name
		ORDER BY count DESC
	`).Scan(&boardRows)

	for _, br := range boardRows {
		result.ContentCountByBoard = append(result.ContentCountByBoard, vo.BoardContentCount{
			BoardID:   br.BoardID,
			BoardName: br.BoardName,
			Count:     br.Count,
		})
	}

	// 质量分分布
	buckets := []struct {
		label string
		lo, hi float64
	}{
		{"0-1", 0, 1},
		{"1-3", 1, 3},
		{"3-5", 3, 5},
		{"5-8", 5, 8},
		{"8+", 8, 9999},
	}
	for _, b := range buckets {
		var cnt int64
		s.db.WithContext(ctx).Raw(
			"SELECT COUNT(*) FROM content_features WHERE quality_score >= ? AND quality_score < ? AND deleted_at IS NULL",
			b.lo, b.hi,
		).Scan(&cnt)
		result.QualityDistribution = append(result.QualityDistribution, vo.ScoreBucket{
			Range: b.label,
			Count: cnt,
		})
	}

	return result, nil
}

// GetRiskAnalysis 获取风控关联分析（结合推荐行为与风险/违规数据）
func (s *recommendationService) GetRiskAnalysis(ctx context.Context) (*vo.RiskAnalysis, error) {
	result := &vo.RiskAnalysis{}

	// 风险用户统计
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM user_risk_records WHERE deleted_at IS NULL",
	).Scan(&result.TotalRiskUsers)

	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM user_risk_records WHERE risk_level = 'danger' AND deleted_at IS NULL",
	).Scan(&result.DangerLevelUsers)

	// 违规统计
	s.db.WithContext(ctx).Model(&do.Violation{}).Count(&result.TotalViolations)
	s.db.WithContext(ctx).Model(&do.Violation{}).
		Where("status = 'pending'").Count(&result.PendingViolations)

	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM violations WHERE punish_type IN ('ban','permanent_ban') AND deleted_at IS NULL",
	).Scan(&result.TotalBans)

	// 违规类型分布
	type violTypeRow struct {
		ViolType string
		Cnt      int64
	}
	var violTypes []violTypeRow
	s.db.WithContext(ctx).Raw(`
		SELECT violation_type as viol_type, COUNT(*) as cnt
		FROM violations WHERE deleted_at IS NULL
		GROUP BY violation_type ORDER BY cnt DESC
	`).Scan(&violTypes)
	var violTotal int64
	for _, vt := range violTypes {
		violTotal += vt.Cnt
	}
	for _, vt := range violTypes {
		ratio := 0.0
		if violTotal > 0 {
			ratio = math.Round(float64(vt.Cnt)/float64(violTotal)*10000) / 100
		}
		result.ViolationDistribution = append(result.ViolationDistribution, vo.ViolationTypeItem{
			ViolationType: vt.ViolType,
			Count:         vt.Cnt,
			Ratio:         ratio,
		})
	}

	// 风控用户行为关联：找出有风险记录的活跃用户
	type riskBehaviorRow struct {
		UserID    uint
		Username  string
		RiskLevel string
		ViolCnt   int64
		BehavCnt  int64
	}
	var riskUsers []riskBehaviorRow
	s.db.WithContext(ctx).Raw(`
		SELECT urr.user_id, COALESCE(u.username, '') as username,
		       urr.risk_level,
		       COALESCE(v.viol_cnt, 0) as viol_cnt,
		       COALESCE(b.behav_cnt, 0) as behav_cnt
		FROM user_risk_records urr
		LEFT JOIN users u ON u.id = urr.user_id AND u.deleted_at IS NULL
		LEFT JOIN (
			SELECT user_id, COUNT(*) as viol_cnt FROM violations WHERE deleted_at IS NULL GROUP BY user_id
		) v ON v.user_id = urr.user_id
		LEFT JOIN (
			SELECT user_id, COUNT(*) as behav_cnt FROM user_behavior_events WHERE deleted_at IS NULL GROUP BY user_id
		) b ON b.user_id = urr.user_id
		WHERE urr.deleted_at IS NULL AND urr.risk_level != 'normal'
		GROUP BY urr.user_id, u.username, urr.risk_level, v.viol_cnt, b.behav_cnt
		ORDER BY COALESCE(v.viol_cnt, 0) DESC, COALESCE(b.behav_cnt, 0) DESC
		LIMIT 10
	`).Scan(&riskUsers)

	userIDs := make([]uint, len(riskUsers))
	for i, ru := range riskUsers {
		userIDs[i] = ru.UserID
	}
	topBehaviorMap := make(map[uint]string)
	for _, uid := range userIDs {
		var bType string
		s.db.WithContext(ctx).Raw(`
			SELECT behavior_type FROM user_behavior_events
			WHERE user_id = ? AND deleted_at IS NULL
			GROUP BY behavior_type ORDER BY COUNT(*) DESC LIMIT 1
		`, uid).Scan(&bType)
		topBehaviorMap[uid] = bType
	}

	for _, ru := range riskUsers {
		result.RiskUserBehaviors = append(result.RiskUserBehaviors, vo.RiskUserBehavior{
			UserID:         ru.UserID,
			Username:       ru.Username,
			RiskLevel:      ru.RiskLevel,
			ViolationCount: ru.ViolCnt,
			BehaviorCount:  ru.BehavCnt,
			TopBehavior:    topBehaviorMap[ru.UserID],
		})
	}

	// 被举报最多的内容
	type reportRow struct {
		CreationID uint
		Title      string
		RptCnt     int64
	}
	var reports []reportRow
	s.db.WithContext(ctx).Raw(`
		SELECT target_id as creation_id, COALESCE(c.title, '') as title, COUNT(*) as rpt_cnt
		FROM reports
		LEFT JOIN creations c ON c.id = reports.target_id AND c.deleted_at IS NULL
		WHERE reports.deleted_at IS NULL AND reports.target_type = 'creation'
		GROUP BY target_id, c.title
		ORDER BY rpt_cnt DESC
		LIMIT 10
	`).Scan(&reports)

	for _, rpt := range reports {
		riskLabel := "low"
		if rpt.RptCnt >= 5 {
			riskLabel = "danger"
		} else if rpt.RptCnt >= 3 {
			riskLabel = "observe"
		}
		result.TopReportedContent = append(result.TopReportedContent, vo.ReportedContentItem{
			CreationID:  rpt.CreationID,
			Title:       rpt.Title,
			ReportCount: rpt.RptCnt,
			RiskLevel:   riskLabel,
		})
	}

	return result, nil
}
