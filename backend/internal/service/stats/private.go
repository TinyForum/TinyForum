package stats

import (
	"context"
	"fmt"
	"sync"
	"time"
	"tiny-forum/internal/model/do"

	"golang.org/x/sync/errgroup"
)

// parseDateRange 将 "YYYY-MM-DD" 字符串解析为完整的日期时间字符串
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

// getBaseInfo 并行获取各维度总量
func (s *statsService) getBaseInfo(ctx context.Context) (*do.StatsInfo, error) {
	var wg sync.WaitGroup
	var info do.StatsInfo
	var err1, err2, err3, err4, err5 error
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
func (s *statsService) getRangeStats(ctx context.Context, startDate, endDate time.Time) (*do.StatsTodayInfo, error) {
	var info do.StatsTodayInfo
	g, ctx := errgroup.WithContext(ctx)

	// 新增用户
	g.Go(func() error {
		count, err := s.userRepo.CountByDateRange(ctx, startDate, endDate)
		if err != nil {
			return fmt.Errorf("统计新增用户: %w", err)
		}
		info.NewUser = count
		return nil
	})

	// 新增文章
	g.Go(func() error {
		count, err := s.postRepo.CountByDateRange(ctx, startDate, endDate)
		if err != nil {
			return fmt.Errorf("统计新增文章: %w", err)
		}
		info.NewArticle = count
		return nil
	})

	// 新增评论
	g.Go(func() error {
		count, err := s.commentRepo.CountByDateRange(ctx, startDate, endDate)
		if err != nil {
			return fmt.Errorf("统计新增评论: %w", err)
		}
		info.NewComment = count
		return nil
	})

	// 新增版块
	g.Go(func() error {
		count, err := s.boardRepo.CountByDateRange(ctx, startDate, endDate)
		if err != nil {
			return fmt.Errorf("统计新增版块: %w", err)
		}
		info.NewBoard = count
		return nil
	})

	// 新增标签
	g.Go(func() error {
		count, err := s.tagRepo.CountByDateRange(ctx, startDate, endDate)
		if err != nil {
			return fmt.Errorf("统计新增标签: %w", err)
		}
		info.NewTag = count
		return nil
	})

	// 活跃用户
	g.Go(func() error {
		count, err := s.userRepo.CountActiveByDateRange(ctx, startDate, endDate)
		if err != nil {
			return fmt.Errorf("统计活跃用户: %w", err)
		}
		info.ActiveUser = count
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return &info, nil
}

// getIllegalInfo 获取违规统计（基于 reports 表）
func (s *statsService) getIllegalInfo(_ context.Context, _, _ time.Time) (*do.StatsIllegalInfo, error) {
	// TODO: 注入 ReportRepository 后按 target_type 分组统计
	return &do.StatsIllegalInfo{}, nil
}

// getActiveUserInfo 获取活跃用户列表及发帖/评论数
func (s *statsService) getActiveUserInfo(ctx context.Context, startDate, endDate time.Time, limit int) (*do.StatsActiveUserInfo, error) {
	users, err := s.userRepo.GetActiveUsersByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	list := make([]*do.ActiveUserDetail, 0, len(users))
	for _, u := range users {
		list = append(list, &do.ActiveUserDetail{
			UserID:       int64(u.ID),
			Username:     u.Username,
			Avatar:       u.Avatar,
			LastActiveAt: time.Now(),
		})
	}
	return &do.StatsActiveUserInfo{
		Total: int64(len(list)),
		List:  list,
	}, nil
}

// getHotArticles 获取热门文章列表
func (s *statsService) getHotArticles(ctx context.Context, startDate, endDate time.Time, limit int) ([]*do.HotArticleItem, error) {
	rows, err := s.postRepo.GetHotArticlesByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	list := make([]*do.HotArticleItem, 0, len(rows))
	for _, a := range rows {
		list = append(list, &do.HotArticleItem{
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
func (s *statsService) getHotBoards(ctx context.Context, startDate, endDate time.Time, limit int) ([]*do.HotBoardItem, error) {
	rows, err := s.boardRepo.GetHotBoardsByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}
	list := make([]*do.HotBoardItem, 0, len(rows))
	for _, b := range rows {
		list = append(list, &do.HotBoardItem{
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
