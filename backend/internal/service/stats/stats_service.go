package stats

import (
	"fmt"
	"time"

	"tiny-forum/internal/repository"
)

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
