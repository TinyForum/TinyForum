package stats

import (
	"context"
	"fmt"
	"time"

	"tiny-forum/internal/model"
)

// GetStatsByDate 获取指定日期的统计数据
func (s *StatsService) GetStatsByDate(ctx context.Context, date time.Time, statsType string) (*model.StatsTodayInfo, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
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
func (s *StatsService) GetTrendStats(ctx context.Context, startDate, endDate time.Time, statsType, intervals string) ([]*model.TrendData, error) {
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
