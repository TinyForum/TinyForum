package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
)

// StatsService 统计服务
type StatsService struct {
	statsRepo   *repository.StatsRepository
	postRepo    repository.PostRepository
	tagRepo     *repository.TagRepository
	boardRepo   *repository.BoardRepository
	userRepo    *repository.UserRepository
	commentRepo *repository.CommentRepository
}

func NewStatsService(
	statsRepo *repository.StatsRepository,
	postRepo repository.PostRepository,
	tagRepo *repository.TagRepository,
	boardRepo *repository.BoardRepository,
	userRepo *repository.UserRepository,
	commentRepo *repository.CommentRepository,
) *StatsService {
	return &StatsService{
		statsRepo:   statsRepo,
		postRepo:    postRepo,
		tagRepo:     tagRepo,
		boardRepo:   boardRepo,
		userRepo:    userRepo,
		commentRepo: commentRepo,
	}
}

// ── 公共方法 ──────────────────────────────────────────────────────────────────

// GetStatsByDate 获取指定日期的统计数据
func (s *StatsService) GetStatsByDate(ctx context.Context, date time.Time, statsType string) (*model.StatsTodayInfo, error) {
	// targetDate, err := time.Parse("2006-01-02", date)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	// }
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// ✅ 方法2：获取当天的结束时间（23:59:59.999）
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())
	todayInfo := &model.StatsTodayInfo{}

	var err error
	switch statsType {
	case "users":
		todayInfo.NewUser, err = s.userRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		return todayInfo, err

	case "posts":
		todayInfo.NewArticle, err = s.postRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		return todayInfo, err

	case "comments":
		todayInfo.NewComment, err = s.commentRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		return todayInfo, err

	default: // "all" or empty
		return s.getRangeStats(ctx, startOfDay, endOfDay)
	}
}

// GetTotalStats 获取指定时间范围的汇总统计数据
func (s *StatsService) GetTotalStats(ctx context.Context, startDate, endDate time.Time, statsType string) (*model.StatsInfoResp, error) {
	fmt.Printf("GetTotalStats: start_date=%s, end_date=%s, stats_type=%s", startDate, endDate, statsType)
	resp := &model.StatsInfoResp{
		StatTime: time.Now(),
	}

	// // startStr, endStr, err := parseDateRange(startDate, endDate)
	// if err != nil {
	// 	return nil, fmt.Errorf("parse time range failed: %w", err)
	// }

	// 基础总量统计（并行）
	baseInfo, err := s.getBaseInfo(ctx)
	if err != nil {
		return nil, err
	}
	resp.BaseInfo = baseInfo

	switch statsType {
	case "users":
		count, err := s.userRepo.CountByDateRange(ctx, startDate, endDate)
		if err != nil {
			return nil, err
		}
		resp.TodayInfo = &model.StatsTodayInfo{NewUser: count}

	case "posts":
		count, err := s.postRepo.CountByDateRange(ctx, startDate, endDate)
		if err != nil {
			return nil, err
		}
		resp.TodayInfo = &model.StatsTodayInfo{NewArticle: count}

	case "comments":
		count, err := s.commentRepo.CountByDateRange(ctx, startDate, endDate)
		if err != nil {
			return nil, err
		}
		resp.TodayInfo = &model.StatsTodayInfo{NewComment: count}

	default: // "all"
		todayInfo, err := s.getRangeStats(ctx, startDate, endDate)
		if err != nil {
			return nil, err
		}
		resp.TodayInfo = todayInfo

		// 以下为可选项，单项失败不影响整体响应
		if illegalInfo, err := s.getIllegalInfo(ctx, startDate, endDate); err == nil {
			resp.IllegalInfo = illegalInfo
		}
		if activeUserInfo, err := s.getActiveUserInfo(ctx, startDate, endDate, 10); err == nil {
			resp.ActiveUserInfo = activeUserInfo
		}
		if hotArticles, err := s.getHotArticles(ctx, startDate, endDate, 10); err == nil {
			resp.HotArticles = hotArticles
		}
		if hotBoards, err := s.getHotBoards(ctx, startDate, endDate, 10); err == nil {
			resp.HotBoards = hotBoards
		}
	}

	return resp, nil
}

// GetTrendStats 获取趋势统计数据（按 day / week / month 粒度）
func (s *StatsService) GetTrendStats(ctx context.Context, startDate, endDate time.Time, statsType, intervals string) ([]*model.TrendData, error) {
	// start, err := time.Parse("2006-01-02", startDate)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid start_date: %w", err)
	// }
	// end, err := time.Parse("2006-01-02", endDate)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid end_date: %w", err)
	// }
	// if end.Before(start) {
	// 	return nil, fmt.Errorf("end_date must not be before start_date")
	// }

	dates := generateDateRange(startDate, endDate, intervals)
	trendData := make([]*model.TrendData, 0, len(dates))

	for _, date := range dates {
		// rangeStart, rangeEnd := dateToRangeStrings(date, interval)

		var count int64
		var err error
		switch statsType {
		case "users":
			count, err = s.userRepo.CountByDateRange(ctx, startDate, endDate)
		case "posts":
			count, err = s.postRepo.CountByDateRange(ctx, startDate, endDate)
		case "comments":
			count, err = s.commentRepo.CountByDateRange(ctx, startDate, endDate)
		default:
			continue
		}
		if err != nil {
			continue // 跳过单个数据点错误，保持其余数据完整
		}

		trendData = append(trendData, &model.TrendData{
			Date:  date.Format("2006-01-02"),
			Count: count,
		})
	}

	return trendData, nil
}

// ── 私有辅助方法 ──────────────────────────────────────────────────────────────

// getBaseInfo 并行获取各维度总量
func (s *StatsService) getBaseInfo(ctx context.Context) (*model.StatsInfo, error) {
	var (
		wg                           sync.WaitGroup
		info                         model.StatsInfo
		err1, err2, err3, err4, err5 error
	)

	wg.Add(5)
	go func() { defer wg.Done(); info.TotalUser, err1 = s.userRepo.Count(ctx) }()
	go func() { defer wg.Done(); info.TotalArticle, err2 = s.postRepo.Count(ctx) }()
	go func() { defer wg.Done(); info.TotalComment, err3 = s.commentRepo.Count(ctx) }()
	go func() { defer wg.Done(); info.TotalBoard, err4 = s.boardRepo.Count(ctx) }()
	go func() { defer wg.Done(); info.TotalTag, err5 = s.tagRepo.Count(ctx) }()
	wg.Wait()

	if err1 != nil {
		return nil, fmt.Errorf("count users: %w", err1)
	}
	if err2 != nil {
		return nil, fmt.Errorf("count posts: %w", err2)
	}
	if err3 != nil {
		return nil, fmt.Errorf("count comments: %w", err3)
	}
	if err4 != nil {
		return nil, fmt.Errorf("count boards: %w", err4)
	}
	if err5 != nil {
		return nil, fmt.Errorf("count tags: %w", err5)
	}

	return &info, nil
}

// getRangeStats 并行获取时间段内各维度增量，单项失败不中断整体
func (s *StatsService) getRangeStats(ctx context.Context, startDate, endDate time.Time) (*model.StatsTodayInfo, error) {
	var wg sync.WaitGroup
	var info model.StatsTodayInfo

	wg.Add(6)
	go func() { defer wg.Done(); info.NewUser, _ = s.userRepo.CountByDateRange(ctx, startDate, endDate) }()
	go func() { defer wg.Done(); info.NewArticle, _ = s.postRepo.CountByDateRange(ctx, startDate, endDate) }()
	go func() { defer wg.Done(); info.NewComment, _ = s.commentRepo.CountByDateRange(ctx, startDate, endDate) }()
	go func() { defer wg.Done(); info.NewBoard, _ = s.boardRepo.CountByDateRange(ctx, startDate, endDate) }()
	go func() { defer wg.Done(); info.NewTag, _ = s.tagRepo.CountByDateRange(ctx, startDate, endDate) }()
	go func() {
		defer wg.Done()
		info.ActiveUser, _ = s.userRepo.CountActiveByDateRange(ctx, startDate, endDate)
	}()
	wg.Wait()

	return &info, nil
}

// getIllegalInfo 获取违规统计（基于 reports 表）
func (s *StatsService) getIllegalInfo(_ context.Context, _, _ time.Time) (*model.StatsIllegalInfo, error) {
	// TODO: 注入 ReportRepository 后按 target_type 分组统计
	return &model.StatsIllegalInfo{}, nil
}

// getActiveUserInfo 获取活跃用户列表及发帖/评论数
func (s *StatsService) getActiveUserInfo(ctx context.Context, startDate, endDate time.Time, limit int) (*model.StatsActiveUserInfo, error) {
	users, err := s.userRepo.GetActiveUsersByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	list := make([]*model.ActiveUserDetail, 0, len(users))
	for _, u := range users {
		list = append(list, &model.ActiveUserDetail{
			UserID:   int64(u.ID),
			Username: u.Username,
			Avatar:   u.Avatar,
			// ArticleCount / CommentCount 需要额外查询，暂留 0
			// LastActiveAt 需要 last_login 或 activity log，暂用当前时间
			LastActiveAt: time.Now(),
		})
	}

	return &model.StatsActiveUserInfo{
		Total: int64(len(list)),
		List:  list,
	}, nil
}

// getHotArticles 获取热门文章列表
func (s *StatsService) getHotArticles(ctx context.Context, startDate, endDate time.Time, limit int) ([]*model.HotArticleItem, error) {
	rows, err := s.postRepo.GetHotArticlesByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	list := make([]*model.HotArticleItem, 0, len(rows))
	for _, a := range rows {
		list = append(list, &model.HotArticleItem{
			ID:           a.ID,
			Title:        a.Title,
			BoardID:      a.BoardID,
			BoardName:    a.BoardName,
			AuthorID:     a.AuthorID,
			AuthorName:   a.AuthorName,
			ViewCount:    a.ViewCount,
			CommentCount: a.CommentCount,
			LikeCount:    a.LikeCount,
			Score:        a.ViewCount + a.CommentCount*10 + a.LikeCount*5,
		})
	}

	return list, nil
}

// getHotBoards 获取热门板块列表
func (s *StatsService) getHotBoards(ctx context.Context, startDate, endDate time.Time, limit int) ([]*model.HotBoardItem, error) {
	rows, err := s.boardRepo.GetHotBoardsByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	list := make([]*model.HotBoardItem, 0, len(rows))
	for _, b := range rows {
		list = append(list, &model.HotBoardItem{
			ID:           b.ID,
			Name:         b.Name,
			Icon:         b.Icon,
			ArticleCount: b.ArticleCount,
			CommentCount: b.CommentCount,
			ActiveUser:   b.ActiveUser,
			Score:        b.ArticleCount*10 + b.CommentCount*2 + b.ActiveUser*5,
		})
	}

	return list, nil
}

// ── 辅助函数 ──────────────────────────────────────────────────────────────────

// parseDateRange 将 "YYYY-MM-DD" 字符串解析为完整的日期时间字符串。
// 两个参数均可为空：startDate 为空时默认为 30 天前，endDate 为空时默认为今天。
func parseDateRange(startDate, endDate string) (string, string, error) {
	now := time.Now()

	var start, end time.Time
	var err error

	if startDate == "" {
		start = now.AddDate(0, 0, -30)
	} else {
		start, err = time.ParseInLocation("2006-01-02", startDate, time.Local)
		if err != nil {
			return "", "", fmt.Errorf("invalid start_date %q: %w", startDate, err)
		}
	}

	if endDate == "" {
		end = now
	} else {
		end, err = time.ParseInLocation("2006-01-02", endDate, time.Local)
		if err != nil {
			return "", "", fmt.Errorf("invalid end_date %q: %w", endDate, err)
		}
	}

	if end.Before(start) {
		return "", "", fmt.Errorf("end_date must not be before start_date")
	}

	return start.Format("2006-01-02 00:00:00"),
		end.Format("2006-01-02 23:59:59"),
		nil
}

// dateToRangeStrings 根据粒度将单个日期扩展为完整的起止时间字符串
func dateToRangeStrings(date time.Time, interval string) (string, string) {
	switch interval {
	case "week":
		weekStart := date.AddDate(0, 0, -int(date.Weekday()))
		weekEnd := weekStart.AddDate(0, 0, 6)
		return weekStart.Format("2006-01-02 00:00:00"), weekEnd.Format("2006-01-02 23:59:59")

	case "month":
		monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
		monthEnd := monthStart.AddDate(0, 1, -1)
		return monthStart.Format("2006-01-02 00:00:00"), monthEnd.Format("2006-01-02 23:59:59")

	default: // "day"
		return date.Format("2006-01-02 00:00:00"), date.Format("2006-01-02 23:59:59")
	}
}

// generateDateRange 按粒度生成日期序列（代表各区间的起始日期）
func generateDateRange(start, end time.Time, interval string) []time.Time {
	var dates []time.Time

	switch interval {
	case "week":
		// 对齐到所在周的周日（Go 的 Weekday: Sunday=0）
		cur := start.AddDate(0, 0, -int(start.Weekday()))
		for !cur.After(end) {
			dates = append(dates, cur)
			cur = cur.AddDate(0, 0, 7)
		}
	case "month":
		cur := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.Local)
		endMonth := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.Local)
		for !cur.After(endMonth) {
			dates = append(dates, cur)
			cur = cur.AddDate(0, 1, 0)
		}
	default: // "day"
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			dates = append(dates, d)
		}
	}

	return dates
}
