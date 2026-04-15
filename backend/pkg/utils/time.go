package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// type TimeTools interface {
// 	// ParseTimeRange 解析时间范围，支持多种格式和相对时间
// 	ParseTimeRange(startExpr, endExpr string) (*TimeRange, error)
// 	// ParseTimeExpression 解析单个时间表达式
// 	ParseTimeExpression(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, error)
// 	// ParseRelativeTime 解析相对时间
// 	ParseRelativeTime(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, bool)
// }

// TimeRange 时间范围
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// type TimeParser struct{}

// ParseTimeRange 解析时间范围，支持多种格式和相对时间
// 支持格式：
// 1. 绝对日期: "2024-01-15", "2024/01/15", "2024.01.15", "2024-01-15 15:30:00"
// 2. 相对时间: "today", "yesterday", "last7days", "last30days", "last90days", "thisweek", "lastweek", "thismonth", "lastmonth", "thisyear"
// 3. 组合: start=2024-01-01&end=2024-01-31 或 start=last7days&end=today
func ParseTimeRange(startExpr, endExpr string) (*TimeRange, error) {
	now := time.Now()
	loc := time.Local
	fmt.Printf("now: %v\n", now)

	// 如果都为空，默认最近30天
	if startExpr == "" && endExpr == "" {
		return &TimeRange{
			Start: now.AddDate(0, 0, -30),
			End:   now,
		}, nil
	}

	var start, end time.Time
	var err error

	// 解析开始时间
	if startExpr != "" {
		start, err = ParseTimeExpression(startExpr, now, loc, false)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %w", err)
		}
	} else {
		// 如果有结束时间但没有开始时间，默认开始时间为结束时间前30天
		end, _ = ParseTimeExpression(endExpr, now, loc, true)
		start = end.AddDate(0, 0, -30)
	}

	// 解析结束时间
	if endExpr != "" {
		end, err = ParseTimeExpression(endExpr, now, loc, true)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %w", err)
		}
	} else {
		// 如果有开始时间但没有结束时间，默认结束时间为当前时间
		end = now
	}

	// 确保开始时间不晚于结束时间
	if start.After(end) {
		start, end = end, start
	}

	return &TimeRange{Start: start, End: end}, nil
}

// parseTimeExpression 解析单个时间表达式
func ParseTimeExpression(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, error) {
	expr = strings.TrimSpace(strings.ToLower(expr))

	// 1. 尝试解析相对时间
	if t, ok := ParseRelativeTime(expr, now, loc, isEnd); ok {
		return t, nil
	}

	// 2. 尝试解析绝对时间（支持多种格式）
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02",
		"2006.01.02",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
	}

	for _, format := range formats {
		if t, err := time.ParseInLocation(format, expr, loc); err == nil {
			if isEnd && !strings.Contains(expr, ":") {
				// 如果是结束日期且没有指定时间，设置为当天结束
				return t.Add(23*time.Hour + 59*time.Minute + 59*time.Second), nil
			}
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported time format: %s", expr)
}

// parseRelativeTime 解析相对时间表达式
func ParseRelativeTime(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, bool) {
	// 获取当天的开始和结束时间
	year, month, day := now.Date()
	todayStart := time.Date(year, month, day, 0, 0, 0, 0, loc)
	todayEnd := time.Date(year, month, day, 23, 59, 59, 999999999, loc)

	switch expr {
	case "today":
		if isEnd {
			return todayEnd, true
		}
		return todayStart, true

	case "yesterday":
		yesterday := now.AddDate(0, 0, -1)
		year, month, day := yesterday.Date()
		if isEnd {
			return time.Date(year, month, day, 23, 59, 59, 999999999, loc), true
		}
		return time.Date(year, month, day, 0, 0, 0, 0, loc), true

	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		year, month, day := tomorrow.Date()
		if isEnd {
			return time.Date(year, month, day, 23, 59, 59, 999999999, loc), true
		}
		return time.Date(year, month, day, 0, 0, 0, 0, loc), true

	case "thisweek":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		weekStart := now.AddDate(0, 0, -(weekday - 1))
		startYear, startMonth, startDay := weekStart.Date()
		weekStart = time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, loc)

		if isEnd {
			weekEnd := weekStart.AddDate(0, 0, 6)
			return time.Date(weekEnd.Year(), weekEnd.Month(), weekEnd.Day(), 23, 59, 59, 999999999, loc), true
		}
		return weekStart, true

	case "lastweek":
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		thisWeekStart := now.AddDate(0, 0, -(weekday - 1))
		lastWeekStart := thisWeekStart.AddDate(0, 0, -7)
		lastWeekEnd := lastWeekStart.AddDate(0, 0, 6)

		startYear, startMonth, startDay := lastWeekStart.Date()
		endYear, endMonth, endDay := lastWeekEnd.Date()

		if isEnd {
			return time.Date(endYear, endMonth, endDay, 23, 59, 59, 999999999, loc), true
		}
		return time.Date(startYear, startMonth, startDay, 0, 0, 0, 0, loc), true

	case "thismonth":
		currentYear, currentMonth, _ := now.Date()
		if isEnd {
			lastDay := time.Date(currentYear, currentMonth+1, 0, 23, 59, 59, 999999999, loc)
			return lastDay, true
		}
		return time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, loc), true

	case "lastmonth":
		lastMonth := now.AddDate(0, -1, 0)
		year, month, _ := lastMonth.Date()
		if isEnd {
			lastDay := time.Date(year, month+1, 0, 23, 59, 59, 999999999, loc)
			return lastDay, true
		}
		return time.Date(year, month, 1, 0, 0, 0, 0, loc), true

	case "thisyear":
		currentYear := now.Year()
		if isEnd {
			return time.Date(currentYear, 12, 31, 23, 59, 59, 999999999, loc), true
		}
		return time.Date(currentYear, 1, 1, 0, 0, 0, 0, loc), true

	case "lastyear":
		lastYear := now.Year() - 1
		if isEnd {
			return time.Date(lastYear, 12, 31, 23, 59, 59, 999999999, loc), true
		}
		return time.Date(lastYear, 1, 1, 0, 0, 0, 0, loc), true
	}

	// 匹配 lastNdays 格式 (如 last7days, last30days)
	reLastNDays := regexp.MustCompile(`^last(\d+)days?$`)
	if matches := reLastNDays.FindStringSubmatch(expr); len(matches) == 2 {
		days, _ := strconv.Atoi(matches[1])
		if isEnd {
			return now, true
		}
		return now.AddDate(0, 0, -days), true
	}

	// 匹配 lastNweeks 格式
	reLastNWeeks := regexp.MustCompile(`^last(\d+)weeks?$`)
	if matches := reLastNWeeks.FindStringSubmatch(expr); len(matches) == 2 {
		weeks, _ := strconv.Atoi(matches[1])
		if isEnd {
			return now, true
		}
		return now.AddDate(0, 0, -weeks*7), true
	}

	// 匹配 lastNmonths 格式
	reLastNMonths := regexp.MustCompile(`^last(\d+)months?$`)
	if matches := reLastNMonths.FindStringSubmatch(expr); len(matches) == 2 {
		months, _ := strconv.Atoi(matches[1])
		if isEnd {
			return now, true
		}
		return now.AddDate(0, -months, 0), true
	}

	return time.Time{}, false
}
