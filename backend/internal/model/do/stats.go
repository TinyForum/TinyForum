package do

import "time"

// StatsInfo 系统基础统计信息
type StatsInfo struct {
	TotalUser    int64 `json:"total_user"`    // 总用户数
	TotalArticle int64 `json:"total_article"` // 总文章数
	TotalComment int64 `json:"total_comment"` // 总评论数
	TotalBoard   int64 `json:"total_board"`   // 总板块数
	TotalTag     int64 `json:"total_tag"`     // 总标签数
}

// StatsTodayInfo 今日统计信息
type StatsTodayInfo struct {
	NewUser    int64 `json:"new_user"`    // 今日新增用户
	NewArticle int64 `json:"new_article"` // 今日新增文章
	NewComment int64 `json:"new_comment"` // 今日新增评论
	NewBoard   int64 `json:"new_board"`   // 今日新增板块
	NewTag     int64 `json:"new_tag"`     // 今日新增标签
	ActiveUser int64 `json:"active_user"` // 今日活跃用户数
}

// StatsInfoResp 统计信息响应（聚合根）
type StatsInfoResp struct {
	BaseInfo       *StatsInfo           `json:"base_info"`                  // 基础统计信息
	TodayInfo      *StatsTodayInfo      `json:"today_info"`                 // 今日统计信息
	IllegalInfo    *StatsIllegalInfo    `json:"illegal_info,omitempty"`     // 今日违规信息
	ActiveUserInfo *StatsActiveUserInfo `json:"active_user_info,omitempty"` // 今日活跃用户信息
	HotArticles    []*HotArticleItem    `json:"hot_articles,omitempty"`     // 今日热门文章列表
	HotBoards      []*HotBoardItem      `json:"hot_boards,omitempty"`       // 今日热门板块列表
	StatTime       time.Time            `json:"stat_time"`                  // 统计时间
}

// StatsIllegalInfo 违规统计信息
type StatsIllegalInfo struct {
	Total      int64 `json:"total"`       // 今日违规总数
	UserCount  int64 `json:"user_count"`  // 今日违规用户数
	ArticleCnt int64 `json:"article_cnt"` // 今日违规文章数
	CommentCnt int64 `json:"comment_cnt"` // 今日违规评论数
	BoardCnt   int64 `json:"board_cnt"`   // 今日违规板块数
}

// StatsActiveUserInfo 活跃用户信息
type StatsActiveUserInfo struct {
	Total int64               `json:"total"` // 今日活跃用户总数
	List  []*ActiveUserDetail `json:"list"`  // 今日活跃用户列表（最多N条）
}

// ActiveUserDetail 活跃用户详情
type ActiveUserDetail struct {
	UserID       int64     `json:"user_id"`        // 用户ID
	Username     string    `json:"username"`       // 用户名
	Avatar       string    `json:"avatar"`         // 头像
	ArticleCount int       `json:"article_count"`  // 今日发文数
	CommentCount int       `json:"comment_count"`  // 今日评论数
	LastActiveAt time.Time `json:"last_active_at"` // 最后活跃时间
}

// HotArticleItem 热门文章项
type HotArticleItem struct {
	ID           int64  `json:"id"`            // 文章ID
	Title        string `json:"title"`         // 文章标题
	BoardID      int64  `json:"board_id"`      // 板块ID
	BoardName    string `json:"board_name"`    // 板块名称
	AuthorID     int64  `json:"author_id"`     // 作者ID
	AuthorName   string `json:"author_name"`   // 作者昵称
	ViewCount    int64  `json:"view_count"`    // 浏览量
	CommentCount int64  `json:"comment_count"` // 评论数
	LikeCount    int64  `json:"like_count"`    // 点赞数
	Score        int64  `json:"score"`         // 综合热度分
}

// HotBoardItem 热门板块项
type HotBoardItem struct {
	ID           int64  `json:"id"`            // 板块ID
	Name         string `json:"name"`          // 板块名称
	Icon         string `json:"icon"`          // 板块图标
	ArticleCount int64  `json:"article_count"` // 今日发文数
	CommentCount int64  `json:"comment_count"` // 今日评论数
	ActiveUser   int64  `json:"active_user"`   // 今日活跃用户数
	Score        int64  `json:"score"`         // 综合热度分
}

// ViolatorItem 违规用户项（可选扩展）
type StatsViolatorItem struct {
	UserID       int64  `json:"user_id"`       // 用户ID
	Username     string `json:"username"`      // 用户名
	ViolationCnt int64  `json:"violation_cnt"` // 违规次数
}

type StatsDayResponse struct {
	Day   string `json:"day"`
	Type  string `json:"type"`
	Count int64  `json:"count"`
}
type StatsTotalResponse struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type StatsTrendResponse struct {
	StartDate string       `json:"start_date"` // 开始日期
	EndDate   string       `json:"end_date"`   // 结束日期
	Interval  string       `json:"interval"`   // 统计粒度 (day/week/month)
	Type      string       `json:"type"`       // 统计类型
	Trend     []*TrendData `json:"trend"`      // 趋势数据
}

// TrendData 趋势数据
type TrendData struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
