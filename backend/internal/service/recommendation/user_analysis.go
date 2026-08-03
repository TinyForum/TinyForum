package recommendation

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/vo"
)

// GetComprehensiveUserAnalysis 获取用户综合分析（画像+标签+行为+相似+风险）
func (s *recommendationService) GetComprehensiveUserAnalysis(ctx context.Context) (*vo.ComprehensiveUserAnalysis, error) {
	result := &vo.ComprehensiveUserAnalysis{}

	now := time.Now()
	since24h := now.Add(-24 * time.Hour)

	// ---- 概览 ----

	// 总用户数
	s.db.WithContext(ctx).Model(&do.User{}).Count(&result.Overview.TotalUsers)

	// 行为记录总数
	s.db.WithContext(ctx).Model(&do.UserBehaviorEvent{}).Count(&result.Overview.TotalBehaviorRecords)

	// 有画像的用户数
	s.db.WithContext(ctx).Model(&do.UserInterestProfile{}).
		Where("tag_weights_json != '{}'").Count(&result.Overview.UsersWithProfile)

	// 人均标签数
	type tagAgg struct{ AvgTags float64 }
	var avgTags tagAgg
	s.db.WithContext(ctx).Raw(`
		SELECT COALESCE(AVG(jsonb_array_length(active_tags_json::jsonb)), 0)
		FROM user_interest_profiles WHERE deleted_at IS NULL AND active_tags_json != '[]'
	`).Scan(&avgTags)
	result.Overview.AvgTagsPerUser = avgTags.AvgTags

	// 24h 人均行为数
	var behav24hCnt int64
	s.db.WithContext(ctx).Model(&do.UserBehaviorEvent{}).
		Where("created_at >= ?", since24h).Count(&behav24hCnt)
	var users24hCnt int64
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE created_at >= ?", since24h,
	).Scan(&users24hCnt)
	if users24hCnt > 0 {
		result.Overview.AvgBehaviorsPerUser24h = math.Round(float64(behav24hCnt)/float64(users24hCnt)*10) / 10
	}

	// 有行为的用户数
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE deleted_at IS NULL",
	).Scan(&result.Overview.SharedUsersCount)

	// 风险用户数
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM user_risk_records WHERE risk_level != 'normal' AND deleted_at IS NULL",
	).Scan(&result.Overview.RiskUserCount)

	// 违规用户数
	s.db.WithContext(ctx).Raw(
		"SELECT COUNT(DISTINCT user_id) FROM violations WHERE deleted_at IS NULL",
	).Scan(&result.Overview.ViolationUserCount)

	// ---- 标签用户分布 ----
	type tagDistRow struct {
		TagID     uint
		TagName   string
		UserCount int64
		AvgWeight float64
		PostCount int64
	}
	var tagDist []tagDistRow
	s.db.WithContext(ctx).Raw(`
		SELECT t.id as tag_id, t.name as tag_name,
		       COALESCE(td.user_cnt, 0) as user_count,
		       COALESCE(td.avg_w, 0) as avg_weight,
		       COALESCE(t.post_count, 0) as post_count
		FROM tags t
		LEFT JOIN (
			SELECT tag_id, COUNT(*) as post_count
			FROM creation_tags
			INNER JOIN creations ON creations.id = creation_tags.creation_id AND creations.deleted_at IS NULL
			WHERE creation_tags.deleted_at IS NULL
			GROUP BY tag_id
		) ctp ON ctp.tag_id = t.id
		LEFT JOIN (
			SELECT tid::int as tag_id, COUNT(user_id) as user_cnt, AVG(w) as avg_w
			FROM (
				SELECT user_id, (jsonb_each_text(tag_weights_json::jsonb)).key as tid,
				       (jsonb_each_text(tag_weights_json::jsonb)).value::float as w
				FROM user_interest_profiles WHERE deleted_at IS NULL AND tag_weights_json != '{}'
			) sub WHERE w > 0
			GROUP BY tid
		) td ON td.tag_id = t.id
		WHERE t.deleted_at IS NULL
		ORDER BY user_count DESC, post_count DESC
		LIMIT 20
	`).Scan(&tagDist)

	for _, td := range tagDist {
		result.TagDistribution = append(result.TagDistribution, vo.TagUserDistribution{
			TagID:     td.TagID,
			TagName:   td.TagName,
			UserCount: td.UserCount,
			AvgWeight: td.AvgWeight,
			PostCount: td.PostCount,
		})
	}

	// ---- 用户行为模式 TOP20 ----
	type behavPatternRow struct {
		UserID          uint
		Username        string
		Nickname        string
		Avatar          string
		TotalBehaviors  int64
		BehavJSON       string
		ActiveTagsJSON  string
		LastActiveAt    time.Time
		RiskLevel       string
		ViolationCount  int64
	}
	var patterns []behavPatternRow
	s.db.WithContext(ctx).Raw(`
		SELECT u.id as user_id, u.username, u.nickname, u.avatar,
		       COUNT(ube.id) as total_behaviors,
		       COALESCE(uip.active_tags_json, '[]') as active_tags_json,
		       MAX(ube.created_at) as last_active_at,
		       COALESCE(urr.risk_level, 'normal') as risk_level,
		       COALESCE(vc.vcnt, 0) as violation_count
		FROM users u
		INNER JOIN user_behavior_events ube ON ube.user_id = u.id AND ube.deleted_at IS NULL
		LEFT JOIN user_interest_profiles uip ON uip.user_id = u.id AND uip.deleted_at IS NULL
		LEFT JOIN user_risk_records urr ON urr.user_id = u.id
			AND urr.risk_level != 'normal' AND urr.deleted_at IS NULL
		LEFT JOIN (
			SELECT user_id, COUNT(*) as vcnt FROM violations WHERE deleted_at IS NULL GROUP BY user_id
		) vc ON vc.user_id = u.id
		WHERE u.deleted_at IS NULL
		GROUP BY u.id, u.username, u.nickname, u.avatar, uip.active_tags_json, urr.risk_level, vc.vcnt
		ORDER BY total_behaviors DESC
		LIMIT 20
	`).Scan(&patterns)

	// 查行为分布 (最多20个用户，一次批量查)
	userIDs := make([]uint, len(patterns))
	for i, p := range patterns {
		userIDs[i] = p.UserID
	}
	behavMap := make(map[uint]map[string]int64)
	if len(userIDs) > 0 {
		type behavItem struct {
			UserID       uint
			BehaviorType string
			Cnt          int64
		}
		var behavItems []behavItem
		s.db.WithContext(ctx).Raw(`
			SELECT user_id, behavior_type, COUNT(*) as cnt
			FROM user_behavior_events
			WHERE user_id IN ? AND deleted_at IS NULL
			GROUP BY user_id, behavior_type
		`, userIDs).Scan(&behavItems)
		for _, bi := range behavItems {
			if behavMap[bi.UserID] == nil {
				behavMap[bi.UserID] = make(map[string]int64)
			}
			behavMap[bi.UserID][bi.BehaviorType] = bi.Cnt
		}
	}

	for _, p := range patterns {
		var activeTags []string
		if p.ActiveTagsJSON != "" {
			json.Unmarshal([]byte(p.ActiveTagsJSON), &activeTags)
		}
		result.UserBehaviorPatterns = append(result.UserBehaviorPatterns, vo.UserBehaviorPattern{
			UserID:           p.UserID,
			Username:         p.Username,
			Nickname:         p.Nickname,
			Avatar:           p.Avatar,
			TotalBehaviors:   p.TotalBehaviors,
			BehaviorBreakdown: behavMap[p.UserID],
			ActiveTags:       activeTags,
			LastActiveAt:     p.LastActiveAt.Format(time.RFC3339),
			RiskLevel:        p.RiskLevel,
			ViolationCount:   p.ViolationCount,
		})
	}

	// ---- 相似用户群组（基于行为的相似度） ----
	// 取行为最活跃的前5个用户，找他们的相似用户
	topActiveUserIDs := userIDs
	if len(topActiveUserIDs) > 5 {
		topActiveUserIDs = topActiveUserIDs[:5]
	}
	for _, uid := range topActiveUserIDs {
		similarIDs, err := s.repo.FindSimilarUsers(uid, 5)
		if err != nil || len(similarIDs) == 0 {
			continue
		}

		type similarDetail struct {
			ID       uint
			Username string
			Nickname string
			Avatar   string
		}
		var details []similarDetail
		s.db.WithContext(ctx).Raw(
			"SELECT id, username, nickname, avatar FROM users WHERE id IN ? AND deleted_at IS NULL", similarIDs,
		).Scan(&details)

		detailMap := make(map[uint]similarDetail)
		for _, d := range details {
			detailMap[d.ID] = d
		}

		// 查种子用户的信息
		var seedUser similarDetail
		s.db.WithContext(ctx).Raw(
			"SELECT id, username, nickname, avatar FROM users WHERE id = ? AND deleted_at IS NULL", uid,
		).Scan(&seedUser)

		var similarItems []vo.SimilarUserItem
		for i, sid := range similarIDs {
			d := detailMap[sid]
			name := d.Nickname
			if name == "" {
				name = d.Username
			}
			similarItems = append(similarItems, vo.SimilarUserItem{
				UserID:          sid,
				Username:        d.Username,
				Nickname:        name,
				Avatar:          d.Avatar,
				SimilarityScore: math.Round(float64(5-i)/float64(5)*100) / 100,
			})
		}

		seedName := seedUser.Nickname
		if seedName == "" {
			seedName = seedUser.Username
		}

		result.SimilarUserGroups = append(result.SimilarUserGroups, vo.SimilarUserGroup{
			SeedUserID:   uid,
			SeedUsername: seedName,
			SimilarUsers: similarItems,
		})
	}

	// ---- 用户风险画像 TOP15 ----
	type riskProfileRow struct {
		UserID            uint
		Username          string
		RiskLevel         string
		ViolationCount    int64
		BehaviorCount     int64
		LastViolationType string
		LastViolationAt   time.Time
		IsBanned          bool
	}
	var riskProfiles []riskProfileRow
	s.db.WithContext(ctx).Raw(`
		SELECT v.user_id, COALESCE(u.username, '') as username,
		       COALESCE(urr.risk_level, 'normal') as risk_level,
		       COUNT(v.id) as violation_count,
		       COALESCE(behav.cnt, 0) as behavior_count,
		       (SELECT violation_type FROM violations WHERE user_id = v.user_id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1) as last_violation_type,
		       COALESCE((SELECT created_at FROM violations WHERE user_id = v.user_id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1), NOW()) as last_violation_at,
		       EXISTS(SELECT 1 FROM violations WHERE user_id = v.user_id AND punish_type IN ('ban','permanent_ban') AND deleted_at IS NULL) as is_banned
		FROM violations v
		LEFT JOIN users u ON u.id = v.user_id AND u.deleted_at IS NULL
		LEFT JOIN user_risk_records urr ON urr.user_id = v.user_id
			AND urr.risk_level != 'normal' AND urr.deleted_at IS NULL
		LEFT JOIN (
			SELECT user_id, COUNT(*) as cnt FROM user_behavior_events WHERE deleted_at IS NULL GROUP BY user_id
		) behav ON behav.user_id = v.user_id
		WHERE v.deleted_at IS NULL
		GROUP BY v.user_id, u.username, urr.risk_level, behav.cnt
		ORDER BY violation_count DESC
		LIMIT 15
	`).Scan(&riskProfiles)

	for _, rp := range riskProfiles {
		lastAt := ""
		if !rp.LastViolationAt.IsZero() {
			lastAt = rp.LastViolationAt.Format(time.RFC3339)
		}
		result.RiskUserList = append(result.RiskUserList, vo.UserRiskProfile{
			UserID:            rp.UserID,
			Username:          rp.Username,
			RiskLevel:         rp.RiskLevel,
			ViolationCount:    rp.ViolationCount,
			BehaviorCount:     rp.BehaviorCount,
			LastViolationType: rp.LastViolationType,
			LastViolationAt:   lastAt,
			IsBanned:          rp.IsBanned,
		})
	}

	return result, nil
}
