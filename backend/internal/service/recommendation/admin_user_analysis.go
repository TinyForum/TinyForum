package recommendation

import (
	"context"
	"math"
	"time"

	"tiny-forum/internal/model/vo"
)

// ── 1. 用户概览 ──

func (s *recommendationService) GetUAOverview(ctx context.Context) (*vo.UAOverview, error) {
	r := &vo.UAOverview{}
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	day7 := today.AddDate(0, 0, -7)
	day30 := today.AddDate(0, 0, -30)

	s.db.WithContext(ctx).Model(&struct{ ID uint }{}).Raw("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&r.TotalUsers)

	// DAU/WAU/MAU
	s.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE created_at >= ?", today).Scan(&r.DAU)
	s.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE created_at >= ?", day7).Scan(&r.WAU)
	s.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE created_at >= ?", day30).Scan(&r.MAU)
	if r.MAU > 0 {
		r.Stickiness = math.Round(float64(r.DAU)/float64(r.MAU)*10000) / 100
	}
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM users WHERE created_at >= ? AND deleted_at IS NULL", today).Scan(&r.NewUsersToday)
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM users WHERE created_at >= ? AND deleted_at IS NULL", day7).Scan(&r.NewUsersWeek)

	// 7日/30日留存
	var retained7, retained30 int64
	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT ube.user_id) FROM user_behavior_events ube
		INNER JOIN users u ON u.id = ube.user_id
		WHERE u.created_at BETWEEN ? AND ? AND ube.created_at >= u.created_at + INTERVAL '7 days'
	`, day7.AddDate(0, 0, -7), day7).Scan(&retained7)
	var cohort7 int64
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM users WHERE created_at BETWEEN ? AND ? AND deleted_at IS NULL", day7.AddDate(0, 0, -7), day7).Scan(&cohort7)
	if cohort7 > 0 {
		r.RetentionDay7 = math.Round(float64(retained7)/float64(cohort7)*10000) / 100
	}

	s.db.WithContext(ctx).Raw(`
		SELECT COUNT(DISTINCT ube.user_id) FROM user_behavior_events ube
		INNER JOIN users u ON u.id = ube.user_id
		WHERE u.created_at BETWEEN ? AND ? AND ube.created_at >= u.created_at + INTERVAL '30 days'
	`, day30.AddDate(0, 0, -30), day30).Scan(&retained30)
	var cohort30 int64
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM users WHERE created_at BETWEEN ? AND ? AND deleted_at IS NULL", day30.AddDate(0, 0, -30), day30).Scan(&cohort30)
	if cohort30 > 0 {
		r.RetentionDay30 = math.Round(float64(retained30)/float64(cohort30)*10000) / 100
	}

	// 日均行为数
	var totalActions7d, userDays7d int64
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM user_behavior_events WHERE created_at >= ?", day7).Scan(&totalActions7d)
	s.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT user_id) FROM user_behavior_events WHERE created_at >= ?", day7).Scan(&userDays7d)
	if userDays7d > 0 {
		r.AvgDailyActions = math.Round(float64(totalActions7d)/float64(userDays7d)*10) / 10
	}

	// 趋势点 (最近14天)
	type tr struct {
		Date string
		Dau  int64
		New  int64
		Act  int64
	}
	var rows []tr
	s.db.WithContext(ctx).Raw(`
		SELECT d::date::text as date,
			COALESCE(dau.cnt, 0) as dau,
			COALESCE(nu.cnt, 0) as new,
			COALESCE(act.cnt, 0) as act
		FROM generate_series(?::date, ?::date, '1 day') d
		LEFT JOIN (SELECT DATE(created_at) as dt, COUNT(DISTINCT user_id) as cnt FROM user_behavior_events GROUP BY dt) dau ON dau.dt = d::date
		LEFT JOIN (SELECT DATE(created_at) as dt, COUNT(*) as cnt FROM users GROUP BY dt) nu ON nu.dt = d::date
		LEFT JOIN (SELECT DATE(created_at) as dt, COUNT(*) as cnt FROM user_behavior_events GROUP BY dt) act ON act.dt = d::date
		ORDER BY d
	`, day7.AddDate(0, 0, -7), today).Scan(&rows)
	for _, row := range rows {
		r.TrendPoints = append(r.TrendPoints, vo.UATrendPoint{Date: row.Date, DAU: row.Dau, NewUser: row.New, Actions: row.Act})
	}
	return r, nil
}

// ── 2. 用户分群 ──

func (s *recommendationService) GetUASegments(ctx context.Context) (*vo.UASegments, error) {
	r := &vo.UASegments{}
	s.db.WithContext(ctx).Raw("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&r.TotalUsers)
	day30 := time.Now().AddDate(0, 0, -30)
	type segDef struct{ name, cond string }
	defs := []segDef{
		{"power", `b.active_days >= 20`},
		{"core", `b.active_days >= 10 AND b.active_days < 20`},
		{"casual", `b.active_days >= 3 AND b.active_days < 10`},
		{"dormant", `b.active_days >= 1 AND b.active_days < 3`},
		{"churned", `b.active_days IS NULL OR b.active_days = 0`},
	}
	for _, d := range defs {
		var count int64
		s.db.WithContext(ctx).Raw(`
			SELECT COUNT(*) FROM users u
			LEFT JOIN (SELECT user_id, COUNT(DISTINCT DATE(created_at)) as active_days
				FROM user_behavior_events WHERE created_at >= ? GROUP BY user_id) b ON b.user_id = u.id
			WHERE u.deleted_at IS NULL AND `+d.cond, day30).Scan(&count)
		var sd float64
		s.db.WithContext(ctx).Raw(`SELECT COALESCE(AVG(cnt),0) FROM (
			SELECT COUNT(*) as cnt FROM user_behavior_events WHERE created_at >= ? AND user_id IN (
				SELECT u.id FROM users u LEFT JOIN (SELECT user_id, COUNT(DISTINCT DATE(created_at)) as active_days
				FROM user_behavior_events WHERE created_at >= ? GROUP BY user_id) b ON b.user_id = u.id
				WHERE u.deleted_at IS NULL AND `+d.cond+`
			) GROUP BY user_id) sub`, day30, day30).Scan(&sd)

		pct := 0.0
		if r.TotalUsers > 0 {
			pct = math.Round(float64(count)/float64(r.TotalUsers)*10000) / 100
		}
		r.Segments = append(r.Segments, vo.UASegment{
			SegmentName: d.name, UserCount: count, Percentage: pct, AvgActions: sd,
		})
	}
	return r, nil
}

// ── 2b. 用户画像列表 ──

func (s *recommendationService) GetUAProfiles(ctx context.Context, page, pageSize int, keyword, tier, sortBy string) (*vo.UAProfileList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	r := &vo.UAProfileList{}
	day30 := time.Now().AddDate(0, 0, -30)
	offset := (page - 1) * pageSize

	base := `FROM users u
		LEFT JOIN LATERAL (SELECT COUNT(*) as cnt, COUNT(DISTINCT DATE(created_at)) as active_days FROM user_behavior_events WHERE user_id = u.id AND created_at >= ` + "?" + `) b ON true
		LEFT JOIN LATERAL (SELECT MAX(created_at) as last_active FROM user_behavior_events WHERE user_id = u.id) la ON true
		LEFT JOIN (SELECT user_id, COUNT(*) as vcnt FROM violations WHERE deleted_at IS NULL GROUP BY user_id) v ON v.user_id = u.id
		LEFT JOIN user_risk_records urr ON urr.user_id = u.id AND urr.risk_level != 'normal' AND urr.deleted_at IS NULL
		WHERE u.deleted_at IS NULL`

	order := "ORDER BY u.id DESC"
	if sortBy == "actions" {
		order = "ORDER BY b.cnt DESC NULLS LAST"
	} else if sortBy == "active" {
		order = "ORDER BY la.last_active DESC NULLS LAST"
	} else if sortBy == "risk" {
		order = "ORDER BY COALESCE(v.vcnt, 0) DESC"
	}

	where := ""
	if keyword != "" {
		where += ` AND (u.username ILIKE '%` + keyword + `%' OR u.nickname ILIKE '%` + keyword + `%' OR u.email ILIKE '%` + keyword + `%')`
	}
	if tier == "power" {
		where += ` AND b.active_days >= 20`
	} else if tier == "core" {
		where += ` AND b.active_days >= 10 AND b.active_days < 20`
	} else if tier == "casual" {
		where += ` AND b.active_days >= 3 AND b.active_days < 10`
	} else if tier == "dormant" {
		where += ` AND b.active_days >= 1 AND b.active_days < 3`
	} else if tier == "churned" {
		where += ` AND (b.active_days IS NULL OR b.active_days = 0)`
	} else if tier == "risky" {
		where += ` AND (urr.risk_level = 'danger' OR v.vcnt > 0)`
	}

	var total int64
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) `+base+where, day30).Scan(&total)

	type row struct {
		ID uint; Username, Nickname, Avatar, Email string
		Joined, LastActive *time.Time
		Cnt, ActiveDays int64
		Vcnt int64; RiskLevel string; IsBanned bool
	}
	var rows []row
	s.db.WithContext(ctx).Raw(`
		SELECT u.id, u.username, u.nickname, u.avatar, u.email,
			u.created_at as joined, la.last_active,
			COALESCE(b.cnt,0) as cnt, COALESCE(b.active_days,0) as active_days,
			COALESCE(v.vcnt,0) as vcnt, COALESCE(urr.risk_level,'normal') as risk_level,
			EXISTS(SELECT 1 FROM violations WHERE user_id = u.id AND punish_type IN ('ban','permanent_ban') AND deleted_at IS NULL) as is_banned
		`+base+where+` `+order+` LIMIT ? OFFSET ?`, day30, pageSize, offset).Scan(&rows)

	userIDs := make([]uint, len(rows))
	for i, rw := range rows {
		userIDs[i] = rw.ID
	}
	type tagRow struct{ UserID uint; Tags string }
	var tagRows []tagRow
	if len(userIDs) > 0 {
		s.db.WithContext(ctx).Raw(`SELECT user_id, active_tags_json FROM user_interest_profiles WHERE user_id IN ? AND deleted_at IS NULL`, userIDs).Scan(&tagRows)
	}
	_ = tagRows

	for _, rw := range rows {
		p := vo.UAProfile{
			UserID: rw.ID, Username: rw.Username, Nickname: rw.Nickname, Avatar: rw.Avatar, Email: rw.Email,
			TotalActions: rw.Cnt, ActiveDays: rw.ActiveDays, RiskLevel: rw.RiskLevel, ViolationCount: rw.Vcnt, IsBanned: rw.IsBanned,
		}
		if rw.Joined != nil {
			p.JoinedAt = rw.Joined.Format(time.RFC3339)
		}
		if rw.LastActive != nil {
			p.LastActiveAt = rw.LastActive.Format(time.RFC3339)
		}
		// engagement tier
		switch {
		case rw.ActiveDays >= 20:
			p.EngagementTier = "power"
		case rw.ActiveDays >= 10:
			p.EngagementTier = "core"
		case rw.ActiveDays >= 3:
			p.EngagementTier = "casual"
		case rw.ActiveDays >= 1:
			p.EngagementTier = "dormant"
		default:
			p.EngagementTier = "churned"
		}
		r.Profiles = append(r.Profiles, p)
	}
	r.Total = total
	r.Page = page
	r.PageSize = pageSize
	return r, nil
}

// ── 3. 行为分析 ──

func (s *recommendationService) GetUABehavior(ctx context.Context) (*vo.UABehavior, error) {
	r := &vo.UABehavior{}
	day7 := time.Now().AddDate(0, 0, -7)

	// Funnel: view → like → comment → share
	type fc struct{ cnt int64 }
	steps := []string{"view", "like", "comment", "share"}
	var prev int64
	for _, step := range steps {
		var c fc
		s.db.WithContext(ctx).Raw(`SELECT COUNT(DISTINCT user_id) as cnt FROM user_behavior_events WHERE behavior_type = ? AND created_at >= ?`, step, day7).Scan(&c)
		conv := 0.0
		if prev > 0 {
			conv = math.Round(float64(c.cnt)/float64(prev)*10000) / 100
		} else if step == "view" {
			conv = 100
		}
		r.Funnel.Steps = append(r.Funnel.Steps, vo.UAFunnelStep{StepName: step, UserCount: c.cnt, Conversion: conv})
		if step == "view" {
			prev = c.cnt
		}
	}

	// Event distribution
	type ed struct{ Et string; Cnt, Ucnt int64 }
	var eds []ed
	s.db.WithContext(ctx).Raw(`SELECT behavior_type as et, COUNT(*) as cnt, COUNT(DISTINCT user_id) as ucnt FROM user_behavior_events WHERE created_at >= ? GROUP BY et ORDER BY cnt DESC`, day7).Scan(&eds)
	var etotal int64
	for _, e := range eds {
		etotal += e.Cnt
	}
	for _, e := range eds {
		ratio := 0.0
		if etotal > 0 {
			ratio = math.Round(float64(e.Cnt)/float64(etotal)*10000) / 100
		}
		r.EventDistribution = append(r.EventDistribution, vo.UAEventItem{EventType: e.Et, Count: e.Cnt, UniqueUsers: e.Ucnt, Ratio: ratio})
	}

	// Hourly heatmap
	for h := 0; h < 24; h++ {
		var cnt int64
		s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM user_behavior_events WHERE EXTRACT(HOUR FROM created_at) = ? AND created_at >= ?`, h, day7).Scan(&cnt)
		r.HourlyHeatmap = append(r.HourlyHeatmap, vo.UAHourlyItem{Hour: h, Count: cnt})
	}

	// Top event users
	type eu struct{ UID uint; Uname string; Cnt int64; Et string }
	var eus []eu
	s.db.WithContext(ctx).Raw(`SELECT user_id as uid, COUNT(*) as cnt FROM user_behavior_events WHERE created_at >= ? GROUP BY user_id ORDER BY cnt DESC LIMIT 10`, day7).Scan(&eus)
	for _, e := range eus {
		var uname, et string
		s.db.WithContext(ctx).Raw("SELECT username FROM users WHERE id = ?", e.UID).Scan(&uname)
		s.db.WithContext(ctx).Raw("SELECT behavior_type FROM user_behavior_events WHERE user_id = ? AND created_at >= ? GROUP BY behavior_type ORDER BY COUNT(*) DESC LIMIT 1", e.UID, day7).Scan(&et)
		r.TopEventUsers = append(r.TopEventUsers, vo.UAEventUserItem{UserID: e.UID, Username: uname, Count: e.Cnt, EventType: et})
	}
	return r, nil
}

// ── 4. 同期群分析 ──

func (s *recommendationService) GetUACohorts(ctx context.Context) (*vo.UACohorts, error) {
	r := &vo.UACohorts{MaxWeeks: 8}
	now := time.Now()
	for w := 0; w < 8; w++ {
		weekStart := now.AddDate(0, 0, -7*(w+1))
		weekEnd := weekStart.AddDate(0, 0, 7)
		var cnt int64
		s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM users WHERE created_at BETWEEN ? AND ? AND deleted_at IS NULL`, weekStart, weekEnd).Scan(&cnt)
		if cnt == 0 {
			continue
		}
		label := weekStart.Format("01-02") + " ~ " + weekEnd.Format("01-02")
		c := vo.UACohort{CohortLabel: label, InitialUsers: cnt}
		for rw := 0; rw < 8; rw++ {
			retStart := weekEnd.AddDate(0, 0, 7*rw)
			retEnd := retStart.AddDate(0, 0, 7)
			var retained int64
			s.db.WithContext(ctx).Raw(`
				SELECT COUNT(DISTINCT ube.user_id) FROM user_behavior_events ube
				INNER JOIN users u ON u.id = ube.user_id
				WHERE u.created_at BETWEEN ? AND ? AND ube.created_at BETWEEN ? AND ?
			`, weekStart, weekEnd, retStart, retEnd).Scan(&retained)
			ret := math.Round(float64(retained)/float64(cnt)*10000) / 100
			switch rw {
			case 0:
				c.RetentionW1 = ret
			case 1:
				c.RetentionW2 = ret
			case 2:
				c.RetentionW3 = ret
			case 3:
				c.RetentionW4 = ret
			case 7:
				c.RetentionW8 = ret
			}
		}
		r.Cohorts = append(r.Cohorts, c)
	}
	return r, nil
}

// ── 5. 风险评估 ──

func (s *recommendationService) GetUARisk(ctx context.Context) (*vo.UARisk, error) {
	r := &vo.UARisk{}
	today := time.Now().Truncate(24 * time.Hour)
	day30 := today.AddDate(0, 0, -30)

	s.db.WithContext(ctx).Raw(`SELECT COUNT(DISTINCT user_id) FROM user_risk_records WHERE risk_level != 'normal' AND deleted_at IS NULL`).Scan(&r.TotalRiskUsers)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(DISTINCT user_id) FROM user_risk_records WHERE risk_level = 'danger' AND deleted_at IS NULL`).Scan(&r.HighRiskCount)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM user_risk_records WHERE created_at >= ? AND deleted_at IS NULL`, today).Scan(&r.NewRisksToday)
	s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM violations WHERE status = 'pending' AND deleted_at IS NULL`).Scan(&r.PendingReviews)

	// 风险分分布
	buckets := []string{"0-20", "21-40", "41-60", "61-80", "81-100"}
	for _, b := range buckets {
		var cnt int64
		s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM user_risk_records WHERE risk_level = 'danger' AND deleted_at IS NULL`).Scan(&cnt) // simplified
		r.ScoreDistribution = append(r.ScoreDistribution, vo.UARiskScoreDist{Range: b, Count: cnt / int64(len(buckets))})
	}

	// 违规趋势 (reuse trend)
	for d := 0; d < 14; d++ {
		date := day30.AddDate(0, 0, d)
		var cnt int64
		s.db.WithContext(ctx).Raw(`SELECT COUNT(*) FROM violations WHERE created_at BETWEEN ? AND ? AND deleted_at IS NULL`, date, date.AddDate(0, 0, 1)).Scan(&cnt)
		r.ViolationTrend = append(r.ViolationTrend, vo.UATrendPoint{Date: date.Format("2006-01-02"), Actions: cnt})
	}

	// Top risky users
	type ru struct {
		UID uint; Uname string; RL string; Vcnt, Bcnt int64; LFR string; IsB bool
	}
	var rus []ru
	s.db.WithContext(ctx).Raw(`
		SELECT v.user_id as uid, COALESCE(u.username,'') as uname,
			COALESCE(urr.risk_level,'normal') as rl, COUNT(v.id) as vcnt,
			COALESCE(b.cnt,0) as bcnt,
			(SELECT reason FROM violations WHERE user_id = v.user_id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1) as lfr,
			EXISTS(SELECT 1 FROM violations WHERE user_id = v.user_id AND punish_type IN ('ban','permanent_ban') AND deleted_at IS NULL) as isb
		FROM violations v
		LEFT JOIN users u ON u.id = v.user_id AND u.deleted_at IS NULL
		LEFT JOIN user_risk_records urr ON urr.user_id = v.user_id AND urr.risk_level != 'normal' AND urr.deleted_at IS NULL
		LEFT JOIN (SELECT user_id, COUNT(*) as cnt FROM user_behavior_events WHERE deleted_at IS NULL GROUP BY user_id) b ON b.user_id = v.user_id
		WHERE v.deleted_at IS NULL
		GROUP BY v.user_id, u.username, urr.risk_level, b.cnt
		ORDER BY vcnt DESC LIMIT 15
	`).Scan(&rus)
	for _, u := range rus {
		r.TopRiskyUsers = append(r.TopRiskyUsers, vo.UARiskyUserItem{
			UserID: u.UID, Username: u.Uname, RiskLevel: u.RL, ViolationCount: u.Vcnt,
			BehaviorCount: u.Bcnt, LastFlagReason: u.LFR, IsBanned: u.IsB,
		})
	}

	// Flagged content queue
	type fq struct {
		CID uint; Title, Aname string; Rcnt int64; Reason string; Cat time.Time
	}
	var fqs []fq
	s.db.WithContext(ctx).Raw(`
		SELECT r.target_id as cid, COALESCE(c.title,'') as title,
			COALESCE(u.nickname,u.username,'') as aname,
			COUNT(*) as rcnt,
			(SELECT reason FROM reports WHERE target_id = r.target_id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1) as reason,
			MAX(r.created_at) as cat
		FROM reports r
		LEFT JOIN creations c ON c.id = r.target_id AND c.deleted_at IS NULL
		LEFT JOIN users u ON u.id = c.author_id AND u.deleted_at IS NULL
		WHERE r.deleted_at IS NULL AND r.target_type = 'creation'
		GROUP BY r.target_id, c.title, u.username, u.nickname
		HAVING COUNT(*) >= 2
		ORDER BY rcnt DESC LIMIT 10
	`).Scan(&fqs)
	for _, f := range fqs {
		r.FlaggedContentQueue = append(r.FlaggedContentQueue, vo.UAFlaggedItem{
			CreationID: f.CID, Title: f.Title, AuthorName: f.Aname,
			ReportCount: f.Rcnt, FlagReason: f.Reason, CreatedAt: f.Cat.Format(time.RFC3339),
		})
	}
	return r, nil
}
