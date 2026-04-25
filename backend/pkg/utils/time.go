package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ============================================================================
// 接口定义（ISP：小而专一）
// ============================================================================

// TimeExpressionParser 解析时间表达式，返回时间和是否成功。
// isEnd 表示该表达式用于结束边界（影响相对时间及无时间的日期处理）。
type TimeExpressionParser interface {
	Parse(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, bool)
}

type TimeHelpers struct {
	SingleParser *TimeParser
	RangeParser  *TimeRangeParser
}

// NewTimeHelpers 创建并返回一个新的TimeHelpers实例
// TimeHelpers是一个时间处理工具结构体，包含单个时间解析器和时间范围解析器
func NewTimeHelpers() *TimeHelpers {
	// 使用默认的解析链（绝对时间 + 相对时间）
	// 创建一个包含绝对时间解析器和相对时间解析器的解析链
	defaultChain := NewParserTimeChain(
		AbsoluteParser{},    // 绝对时间解析器
		NewRelativeParser(), // 相对时间解析器
	)
	// 返回一个新的TimeHelpers实例，初始化其SingleParser和RangeParser字段
	return &TimeHelpers{
		SingleParser: NewTimeParser(defaultChain),      // 单个时间解析器
		RangeParser:  NewTimeRangeParser(defaultChain), // 时间范围解析器
	}
}

// ============================================================================
// 具体解析器实现（SRP：每个解析器只负责一种格式）
// ============================================================================

// AbsoluteParser 解析绝对日期时间字符串。
type AbsoluteParser struct{}

type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Parse 实现 TimeExpressionParser 接口。
func (AbsoluteParser) Parse(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, bool) {
	expr = strings.TrimSpace(expr)
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
				// 结束日期无时间 → 当天 23:59:59
				t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			}
			return t, true
		}
	}
	return time.Time{}, false
}

// RelativeParser 解析相对时间表达式（today, yesterday, last7days, thisweek 等）。
type RelativeParser struct {
	reLastNDays   *regexp.Regexp
	reLastNWeeks  *regexp.Regexp
	reLastNMonths *regexp.Regexp
}

// NewRelativeParser 创建相对时间解析器并预编译正则。
func NewRelativeParser() *RelativeParser {
	return &RelativeParser{
		reLastNDays:   regexp.MustCompile(`^last(\d+)days?$`),
		reLastNWeeks:  regexp.MustCompile(`^last(\d+)weeks?$`),
		reLastNMonths: regexp.MustCompile(`^last(\d+)months?$`),
	}
}

// Parse 实现 TimeExpressionParser 接口。
func (p *RelativeParser) Parse(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, bool) {
	expr = strings.ToLower(strings.TrimSpace(expr))

	// 辅助函数：当天开始/结束
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
		yYear, yMonth, yDay := yesterday.Date()
		if isEnd {
			return time.Date(yYear, yMonth, yDay, 23, 59, 59, 999999999, loc), true
		}
		return time.Date(yYear, yMonth, yDay, 0, 0, 0, 0, loc), true

	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		tYear, tMonth, tDay := tomorrow.Date()
		if isEnd {
			return time.Date(tYear, tMonth, tDay, 23, 59, 59, 999999999, loc), true
		}
		return time.Date(tYear, tMonth, tDay, 0, 0, 0, 0, loc), true

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

	// lastNdays
	if matches := p.reLastNDays.FindStringSubmatch(expr); len(matches) == 2 {
		days, _ := strconv.Atoi(matches[1])
		if isEnd {
			return now, true
		}
		return now.AddDate(0, 0, -days), true
	}

	// lastNweeks
	if matches := p.reLastNWeeks.FindStringSubmatch(expr); len(matches) == 2 {
		weeks, _ := strconv.Atoi(matches[1])
		if isEnd {
			return now, true
		}
		return now.AddDate(0, 0, -weeks*7), true
	}

	// lastNmonths
	if matches := p.reLastNMonths.FindStringSubmatch(expr); len(matches) == 2 {
		months, _ := strconv.Atoi(matches[1])
		if isEnd {
			return now, true
		}
		return now.AddDate(0, -months, 0), true
	}

	return time.Time{}, false
}

// ============================================================================
// 组合解析器（职责链，OCP：可无限扩展）
// ============================================================================

// ParserTimeChain 实现 TimeExpressionParser 接口，按顺序尝试多个解析器。
type ParserTimeChain struct {
	parsers []TimeExpressionParser
}

// NewParserTimeChain 创建一个解析器链。
func NewParserTimeChain(parsers ...TimeExpressionParser) *ParserTimeChain {
	return &ParserTimeChain{parsers: parsers}
}

// Add 在链末尾添加解析器（返回自身，支持链式调用）。
func (c *ParserTimeChain) Add(p TimeExpressionParser) *ParserTimeChain {
	c.parsers = append(c.parsers, p)
	return c
}

// Parse 实现 TimeExpressionParser 接口。
func (c *ParserTimeChain) Parse(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, bool) {
	for _, p := range c.parsers {
		if t, ok := p.Parse(expr, now, loc, isEnd); ok {
			return t, true
		}
	}
	return time.Time{}, false
}

// ============================================================================
// 带错误返回的封装器（方便业务层使用）
// ============================================================================

// TimeParser 封装一个 TimeExpressionParser，提供返回 error 的 Parse 方法。
type TimeParser struct {
	parser TimeExpressionParser
}

// NewTimeParser 创建带错误返回的解析器。
func NewTimeParser(parser TimeExpressionParser) *TimeParser {
	return &TimeParser{parser: parser}
}

// Parse 解析表达式，返回时间或 error。
func (tp *TimeParser) Parse(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, error) {
	if t, ok := tp.parser.Parse(expr, now, loc, isEnd); ok {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %s", expr)
}

// ============================================================================
// 时间范围解析器（DIP：依赖 TimeExpressionParser 接口）
// ============================================================================

// TimeRangeParser 解析时间范围，支持默认值、自动交换等。
type TimeRangeParser struct {
	exprParser TimeExpressionParser
}

// NewTimeRangeParser 创建时间范围解析器。
func NewTimeRangeParser(parser TimeExpressionParser) *TimeRangeParser {
	return &TimeRangeParser{exprParser: parser}
}

// Parse 解析开始和结束表达式，返回 TimeRange。
func (p *TimeRangeParser) Parse(startExpr, endExpr string) (*TimeRange, error) {
	now := time.Now()
	loc := time.Local

	// 都为空 → 最近30天
	if startExpr == "" && endExpr == "" {
		return &TimeRange{
			Start: now.AddDate(0, 0, -30),
			End:   now,
		}, nil
	}

	var start, end time.Time
	var err error

	// 开始时间
	if startExpr != "" {
		start, err = p.parseSingle(startExpr, now, loc, false)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %w", err)
		}
	} else {
		end, err = p.parseSingle(endExpr, now, loc, true)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %w", err)
		}
		start = end.AddDate(0, 0, -30)
	}

	// 结束时间
	if endExpr != "" {
		end, err = p.parseSingle(endExpr, now, loc, true)
		if err != nil {
			return nil, fmt.Errorf("invalid end time: %w", err)
		}
	} else {
		end = now
	}

	// 保证时间顺序
	if start.After(end) {
		start, end = end, start
	}

	return &TimeRange{Start: start, End: end}, nil
}

// parseSingle 辅助方法：将 bool 结果转换为 error。
func (p *TimeRangeParser) parseSingle(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, error) {
	if t, ok := p.exprParser.Parse(expr, now, loc, isEnd); ok {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unsupported time format: %s", expr)
}

// ============================================================================
// 包级默认实例（向后兼容 + 开闭原则：可替换）
// ============================================================================

// var (
// 	// defaultChain 默认解析链：绝对时间 → 相对时间
// 	defaultChain = NewParserTimeChain(AbsoluteParser{}, NewRelativeParser())

// 	// DefaultTimeParser 默认的错误返回解析器
// 	DefaultTimeParser = NewTimeParser(defaultChain)

// 	// DefaultTimeRangeParser 默认的时间范围解析器
// 	DefaultTimeRangeParser = NewTimeRangeParser(defaultChain)
// )

// // ParseTimeExpression 解析单个时间表达式（与旧版 API 完全兼容）。
// func ParseTimeExpression(expr string, now time.Time, loc *time.Location, isEnd bool) (time.Time, error) {
// 	return DefaultTimeParser.Parse(expr, now, loc, isEnd)
// }

// // ParseTimeRange 解析时间范围（与旧版 API 完全兼容）。
// func ParseTimeRange(startExpr, endExpr string) (*TimeRange, error) {
// 	return DefaultTimeRangeParser.Parse(startExpr, endExpr)
// }
