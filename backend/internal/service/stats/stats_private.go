package stats

import (
	"context"
	"fmt"
	"sync"
	"time"

	"tiny-forum/internal/model"
)

// getBaseInfo 并行获取各维度总量
func (s *StatsService) getBaseInfo(ctx context.Context) (*model.StatsInfo, error) {
	var wg sync.WaitGroup
	var info model.StatsInfo
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
			UserID:       int64(u.ID),
			Username:     u.Username,
			Avatar:       u.Avatar,
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
