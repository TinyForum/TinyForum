package stats

import (
	"context"
	"fmt"
	"time"

	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
)

// GetStatsByDate 获取指定日期的统计数据
func (s *statsService) GetStatsByDate(ctx context.Context, date time.Time, statsType string) (*model.StatsTodayInfo, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, date.Location())

	// 统一返回完整结构
	info := &model.StatsTodayInfo{}

	var err error
	switch statsType {
	case "users":
		info.NewUser, err = s.userRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		if err != nil {
			return nil, fmt.Errorf("查询新增用户失败: %w", err)
		}
		return info, nil
	case "posts":
		info.NewArticle, err = s.postRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		if err != nil {
			return nil, fmt.Errorf("查询新增文章失败: %w", err)
		}
		return info, nil
	case "comments":
		info.NewComment, err = s.commentRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		if err != nil {
			return nil, fmt.Errorf("查询新增评论失败: %w", err)
		}
		return info, nil
	default: // "all" or empty
		return s.getRangeStats(ctx, startOfDay, endOfDay)
	}
}

// GetTotalStats 获取指定时间范围的汇总统计数据
func (s *statsService) GetTotalStats(ctx context.Context, startDate, endDate time.Time, statsType string) (*model.StatsInfoResp, error) {
	fmt.Printf("GetTotalStats: start_date=%s, end_date=%s, stats_type=%s", startDate, endDate, statsType)
	resp := &model.StatsInfoResp{StatTime: time.Now()}
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
func (s *statsService) GetTrendStats(ctx context.Context, startDate, endDate time.Time, statsType, intervals string) ([]*model.TrendData, error) {
	dates := generateDateRange(startDate, endDate, intervals)
	trendData := make([]*model.TrendData, 0, len(dates))
	for _, date := range dates {
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
			continue
		}
		trendData = append(trendData, &model.TrendData{
			Date:  date.Format("2006-01-02"),
			Count: count,
		})
	}
	return trendData, nil
}

// StatsService 新增方法
func (s *statsService) GetStatsByDateRange(ctx context.Context, startDate, endDate time.Time, statsType string) ([]dto.DailyStatResponse, error) {
	// 将日期对齐到零点（避免时区漂移）
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	if start.After(end) {
		return nil, fmt.Errorf("开始日期不能晚于结束日期")
	}

	days := int(end.Sub(start).Hours()/24) + 1
	result := make([]dto.DailyStatResponse, 0, days)

	// 逐日统计（可考虑并发优化，但30天内简单循环足够）
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		dailyInfo, err := s.GetStatsByDate(ctx, d, statsType)
		if err != nil {
			return nil, fmt.Errorf("统计日期 %s 失败: %w", d.Format("2006-01-02"), err)
		}

		result = append(result, dto.DailyStatResponse{
			Date:       d.Format("2006-01-02"),
			NewUser:    dailyInfo.NewUser,
			NewArticle: dailyInfo.NewArticle,
			NewComment: dailyInfo.NewComment,
			NewBoard:   dailyInfo.NewBoard,
			NewTag:     dailyInfo.NewTag,
			ActiveUser: dailyInfo.ActiveUser,
		})
	}

	return result, nil
}
