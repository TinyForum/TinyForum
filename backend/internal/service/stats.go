package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
)

type StatsService struct {
	postRepo        repository.PostRepository
	tagRepo         *repository.TagRepository
	boardRepo       *repository.BoardRepository
	userRepo        *repository.UserRepository
	notifSvc        *repository.NotificationRepository
	timelineSvc     *repository.TimelineRepository
	topicSvc        *repository.TopicRepository
	commentSvc      *repository.CommentRepository
	announcementSvc repository.AnnouncementRepository
	questionSvc     repository.QuestionRepository
}

func NewStatsService(
	statsRepo *repository.StatsRepository,
	postRepo repository.PostRepository,
	tagRepo *repository.TagRepository,
	boardRepo *repository.BoardRepository,
	userRepo *repository.UserRepository,
	timelineRepo *repository.TimelineRepository,
	notifRepo *repository.NotificationRepository,
	topicRepo *repository.TopicRepository,
	commentRepo *repository.CommentRepository,
	announcementRepo repository.AnnouncementRepository,
	questionRepo repository.QuestionRepository,

) *StatsService {
	return &StatsService{
		postRepo:        postRepo,
		tagRepo:         tagRepo,
		boardRepo:       boardRepo,
		userRepo:        userRepo,
		notifSvc:        notifRepo,
		timelineSvc:     timelineRepo,
		topicSvc:        topicRepo,
		commentSvc:      commentRepo,
		announcementSvc: announcementRepo,
		questionSvc:     questionRepo,
	}
}

// GetStatsByDate 获取指定日期的统计数据
func (s *StatsService) GetStatsByDate(ctx context.Context, date string, statsType string) (*model.StatsTodayInfo, error) {
	// 解析日期
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, err
	}

	// 计算日期范围（当天 00:00:00 到 23:59:59）
	startOfDay := targetDate.Format("2006-01-02 00:00:00")
	endOfDay := targetDate.Format("2006-01-02 23:59:59")

	todayInfo := &model.StatsTodayInfo{}

	// 根据类型获取不同维度的数据
	switch statsType {
	case "users":
		newUser, err := s.userRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		if err != nil {
			return nil, err
		}
		todayInfo.NewUser = newUser

	case "posts":
		newArticle, err := s.postRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		if err != nil {
			return nil, err
		}
		todayInfo.NewArticle = newArticle

	case "comments":
		newComment, err := s.commentSvc.CountByDateRange(ctx, startOfDay, endOfDay)
		if err != nil {
			return nil, err
		}
		todayInfo.NewComment = newComment

	case "likes":
		// 点赞统计（如果有 LikeRepository）
		// likeCount, err := s.likeRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		// todayInfo.LikeCount = likeCount

	default: // "all"
		// 并行获取所有数据
		var wg sync.WaitGroup
		var newUser, newArticle, newComment, newBoard, newTag, activeUser int64
		var err1, err2, err3, err4, err5, err6 error

		wg.Add(6)

		go func() {
			defer wg.Done()
			newUser, err1 = s.userRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		}()

		go func() {
			defer wg.Done()
			newArticle, err2 = s.postRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		}()

		go func() {
			defer wg.Done()
			newComment, err3 = s.commentSvc.CountByDateRange(ctx, startOfDay, endOfDay)
		}()

		go func() {
			defer wg.Done()
			newBoard, err4 = s.boardRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		}()

		go func() {
			defer wg.Done()
			newTag, err5 = s.tagRepo.CountByDateRange(ctx, startOfDay, endOfDay)
		}()

		go func() {
			defer wg.Done()
			activeUser, err6 = s.userRepo.CountActiveByDateRange(ctx, startOfDay, endOfDay)
		}()

		wg.Wait()

		// 检查错误（可选：部分失败不影响整体）
		if err1 == nil {
			todayInfo.NewUser = newUser
		}
		if err2 == nil {
			todayInfo.NewArticle = newArticle
		}
		if err3 == nil {
			todayInfo.NewComment = newComment
		}
		if err4 == nil {
			todayInfo.NewBoard = newBoard
		}
		if err5 == nil {
			todayInfo.NewTag = newTag
		}
		if err6 == nil {
			todayInfo.ActiveUser = activeUser
		}
	}

	return todayInfo, nil
}

// GetTotalStats 获取总计统计数据（聚合根）
func (s *StatsService) GetTotalStats(ctx context.Context, startDate, endDate string, statsType string) (*model.StatsInfoResp, error) {
	resp := &model.StatsInfoResp{
		StatTime: time.Now(),
	}

	// 解析日期范围
	start, end, err := parseDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 获取基础统计信息（总数）
	baseInfo, err := s.getBaseInfo(ctx)
	if err != nil {
		return nil, err
	}
	resp.BaseInfo = baseInfo

	// 根据类型获取区间统计
	switch statsType {
	case "users":
		// 获取区间内新增用户
		newUser, err := s.userRepo.CountByDateRange(ctx, start, end)
		if err != nil {
			return nil, err
		}
		resp.TodayInfo = &model.StatsTodayInfo{NewUser: newUser}

	case "posts":
		newArticle, err := s.postRepo.CountByDateRange(ctx, start, end)
		if err != nil {
			return nil, err
		}
		resp.TodayInfo = &model.StatsTodayInfo{NewArticle: newArticle}

	case "comments":
		newComment, err := s.commentSvc.CountByDateRange(ctx, start, end)
		if err != nil {
			return nil, err
		}
		resp.TodayInfo = &model.StatsTodayInfo{NewComment: newComment}

	default: // "all"
		// 获取区间内所有统计数据
		todayInfo, err := s.getRangeStats(ctx, start, end)
		if err != nil {
			return nil, err
		}
		resp.TodayInfo = todayInfo

		// 获取违规信息（可选）
		illegalInfo, err := s.getIllegalInfo(ctx, start, end)
		if err == nil {
			resp.IllegalInfo = illegalInfo
		}

		// 获取活跃用户信息
		activeUserInfo, err := s.getActiveUserInfo(ctx, start, end, 10)
		if err == nil {
			resp.ActiveUserInfo = activeUserInfo
		}

		// 获取热门文章
		hotArticles, err := s.getHotArticles(ctx, start, end, 10)
		if err == nil {
			resp.HotArticles = hotArticles
		}

		// 获取热门板块
		hotBoards, err := s.getHotBoards(ctx, start, end, 10)
		if err == nil {
			resp.HotBoards = hotBoards
		}
	}

	return resp, nil
}

// GetTrendStats 获取趋势统计数据
func (s *StatsService) GetTrendStats(ctx context.Context, startDate, endDate, statsType, interval string) ([]*model.TrendData, error) {
	// 解析日期范围
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	// 根据粒度生成时间序列
	dates := generateDateRange(start, end, interval)

	var trendData []*model.TrendData

	for _, date := range dates {
		var dateRangeStart, dateRangeEnd string

		switch interval {
		case "day":
			dateRangeStart = date.Format("2006-01-02 00:00:00")
			dateRangeEnd = date.Format("2006-01-02 23:59:59")
		case "week":
			weekStart := date.AddDate(0, 0, -int(date.Weekday()))
			weekEnd := weekStart.AddDate(0, 0, 6)
			dateRangeStart = weekStart.Format("2006-01-02 00:00:00")
			dateRangeEnd = weekEnd.Format("2006-01-02 23:59:59")
		case "month":
			monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
			monthEnd := monthStart.AddDate(0, 1, -1)
			dateRangeStart = monthStart.Format("2006-01-02 00:00:00")
			dateRangeEnd = monthEnd.Format("2006-01-02 23:59:59")
		}

		var count int64

		switch statsType {
		case "users":
			count, err = s.userRepo.CountByDateRange(ctx, dateRangeStart, dateRangeEnd)
		case "posts":
			count, err = s.postRepo.CountByDateRange(ctx, dateRangeStart, dateRangeEnd)
		case "comments":
			count, err = s.commentSvc.CountByDateRange(ctx, dateRangeStart, dateRangeEnd)
		case "likes":
			// count, err = s.likeRepo.CountByDateRange(ctx, dateRangeStart, dateRangeEnd)
			count = 0
		}

		if err != nil {
			continue // 跳过错误的数据点
		}

		trendData = append(trendData, &model.TrendData{
			Date:  date.Format("2006-01-02"),
			Count: count,
		})
	}

	return trendData, nil
}

// ==================== 私有辅助方法 ====================

// getBaseInfo 获取基础统计信息（总数）
func (s *StatsService) getBaseInfo(ctx context.Context) (*model.StatsInfo, error) {
	var wg sync.WaitGroup
	var info model.StatsInfo
	var err1, err2, err3, err4, err5 error

	wg.Add(5)

	go func() {
		defer wg.Done()
		info.TotalUser, err1 = s.userRepo.Count(ctx)
	}()

	go func() {
		defer wg.Done()
		info.TotalArticle, err2 = s.postRepo.Count(ctx)
	}()

	go func() {
		defer wg.Done()
		info.TotalComment, err3 = s.commentSvc.Count(ctx)
	}()

	go func() {
		defer wg.Done()
		info.TotalBoard, err4 = s.boardRepo.Count(ctx)
	}()

	go func() {
		defer wg.Done()
		info.TotalTag, err5 = s.tagRepo.Count(ctx)
	}()

	wg.Wait()

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		return nil, fmt.Errorf("获取基础统计失败")
	}

	return &info, nil
}

// getRangeStats 获取区间统计数据
func (s *StatsService) getRangeStats(ctx context.Context, startDate, endDate string) (*model.StatsTodayInfo, error) {
	var wg sync.WaitGroup
	var info model.StatsTodayInfo

	wg.Add(6)

	go func() {
		defer wg.Done()
		info.NewUser, _ = s.userRepo.CountByDateRange(ctx, startDate, endDate)
	}()

	go func() {
		defer wg.Done()
		info.NewArticle, _ = s.postRepo.CountByDateRange(ctx, startDate, endDate)
	}()

	go func() {
		defer wg.Done()
		info.NewComment, _ = s.commentSvc.CountByDateRange(ctx, startDate, endDate)
	}()

	go func() {
		defer wg.Done()
		info.NewBoard, _ = s.boardRepo.CountByDateRange(ctx, startDate, endDate)
	}()

	go func() {
		defer wg.Done()
		info.NewTag, _ = s.tagRepo.CountByDateRange(ctx, startDate, endDate)
	}()

	go func() {
		defer wg.Done()
		info.ActiveUser, _ = s.userRepo.CountActiveByDateRange(ctx, startDate, endDate)
	}()

	wg.Wait()

	return &info, nil
}

// getIllegalInfo 获取违规统计信息
func (s *StatsService) getIllegalInfo(ctx context.Context, startDate, endDate string) (*model.StatsIllegalInfo, error) {
	// 需要根据实际的违规记录表来实现
	// 这里假设有对应的 repository 方法

	info := &model.StatsIllegalInfo{}

	// 示例实现（需要根据实际表结构调整）
	// total, err := s.reportRepo.CountByDateRange(ctx, startDate, endDate)
	// if err != nil {
	//     return nil, err
	// }
	// info.Total = total

	return info, nil
}

// getActiveUserInfo 获取活跃用户信息
// getActiveUserInfo 获取活跃用户信息
func (s *StatsService) getActiveUserInfo(ctx context.Context, startDate, endDate string, limit int) (*model.StatsActiveUserInfo, error) {
	// 获取活跃用户列表
	users, err := s.userRepo.GetActiveUsersByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	var list []*model.ActiveUserDetail
	for _, u := range users {
		list = append(list, &model.ActiveUserDetail{
			UserID:       int64(u.ID), // 修复：uint -> int64 类型转换
			Username:     u.Username,
			Avatar:       u.Avatar,
			ArticleCount: 0,          // User 模型没有 PostCount 字段，暂时设为 0
			CommentCount: 0,          // User 模型没有 CommentCount 字段，暂时设为 0
			LastActiveAt: time.Now(), // User 模型没有 LastActiveAt 字段，使用当前时间
		})
	}

	return &model.StatsActiveUserInfo{
		Total: int64(len(list)),
		List:  list,
	}, nil
}

// getHotArticles 获取热门文章
func (s *StatsService) getHotArticles(ctx context.Context, startDate, endDate string, limit int) ([]*model.HotArticleItem, error) {
	// 获取热门文章（按综合热度排序）
	articles, err := s.postRepo.GetHotArticlesByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	var list []*model.HotArticleItem
	for _, a := range articles {
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
			Score:        a.ViewCount + a.CommentCount*10 + a.LikeCount*5, // 综合热度分算法
		})
	}

	return list, nil
}

// getHotBoards 获取热门板块
// getHotBoards 获取热门板块
func (s *StatsService) getHotBoards(ctx context.Context, startDate, endDate string, limit int) ([]*model.HotBoardItem, error) {
	// 获取热门板块（按活跃度排序）
	boards, err := s.boardRepo.GetHotBoardsByDateRange(ctx, startDate, endDate, limit)
	if err != nil {
		return nil, err
	}

	var list []*model.HotBoardItem
	for _, b := range boards {
		list = append(list, &model.HotBoardItem{
			ID:           b.ID,
			Name:         b.Name,
			Icon:         b.Icon,
			ArticleCount: b.ArticleCount, // 修复：使用 ArticleCount 而不是 PostCount
			CommentCount: b.CommentCount,
			ActiveUser:   b.ActiveUser, // 修复：使用 ActiveUser 而不是 ActiveUserCount
			Score:        b.ArticleCount*10 + b.CommentCount*2 + b.ActiveUser*5,
		})
	}

	return list, nil
}

// ==================== 辅助函数 ====================

func parseDateRange(startDate, endDate string) (string, string, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return "", "", err
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return "", "", err
	}

	startStr := start.Format("2006-01-02 00:00:00")
	endStr := end.Format("2006-01-02 23:59:59")

	return startStr, endStr, nil
}

func generateDateRange(start, end time.Time, interval string) []time.Time {
	var dates []time.Time

	switch interval {
	case "day":
		for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
			dates = append(dates, d)
		}
	case "week":
		// 获取起始周的周一
		startWeek := start.AddDate(0, 0, -int(start.Weekday()))
		for d := startWeek; !d.After(end); d = d.AddDate(0, 0, 7) {
			dates = append(dates, d)
		}
	case "month":
		current := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.Local)
		endMonth := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.Local)
		for !current.After(endMonth) {
			dates = append(dates, current)
			current = current.AddDate(0, 1, 0)
		}
	}

	return dates
}
