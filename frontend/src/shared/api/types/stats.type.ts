// stats.ts

/**
 * 系统基础统计信息
 */
export interface StatsInfo {
  total_user: number;
  total_article: number;
  total_comment: number;
  total_board: number;
  total_tag: number;
}

/**
 * 今日统计信息
 */
export interface StatsTodayInfo {
  new_user: number;
  new_article: number;
  new_comment: number;
  new_board: number;
  new_tag: number;
  active_user: number;
}

/**
 * 违规统计信息
 */
export interface StatsIllegalInfo {
  total: number;
  user_count: number;
  article_cnt: number;
  comment_cnt: number;
  board_cnt: number;
}

/**
 * 活跃用户详情
 */
export interface ActiveUserDetail {
  user_id: number;
  username: string;
  avatar: string;
  article_count: number;
  comment_count: number;
  last_active_at: string; // ISO 8601 时间字符串
}

/**
 * 活跃用户信息
 */
export interface StatsActiveUserInfo {
  total: number;
  list: ActiveUserDetail[];
}

/**
 * 热门文章项
 */
export interface HotArticleItem {
  id: number;
  title: string;
  board_id: number;
  board_name: string;
  author_id: number;
  author_name: string;
  view_count: number;
  comment_count: number;
  like_count: number;
  score: number;
}

/**
 * 热门板块项
 */
export interface HotBoardItem {
  id: number;
  name: string;
  icon: string;
  article_count: number;
  comment_count: number;
  active_user: number;
  score: number;
}

/**
 * 违规用户项
 */
export interface StatsViolatorItem {
  user_id: number;
  username: string;
  violation_cnt: number;
}

/**
 * 统计信息响应（聚合根）
 */
export interface StatsInfoResp {
  base_info: StatsInfo;
  today_info: StatsTodayInfo;
  illegal_info?: StatsIllegalInfo;
  active_user_info?: StatsActiveUserInfo;
  hot_articles?: HotArticleItem[];
  hot_boards?: HotBoardItem[];
  stat_time: string; // ISO 8601 时间字符串
}

/**
 * 日统计数据响应
 */
export interface StatsDayResponse {
  day: string;
  type: string;
  count: number;
}

/**
 * 总计统计数据响应
 */
export interface StatsTotalResponse {
  type: string;
  count: number;
}

/**
 * 趋势数据
 */
export interface TrendData {
  date: string;
  count: number;
}

/**
 * 趋势统计响应
 */
export interface StatsTrendResponse {
  start_date: string;
  end_date: string;
  interval: "day" | "week" | "month";
  type: "users" | "posts" | "comments" | "likes";
  trend: TrendData[];
}

/**
 * API 统一响应格式
 */
// export interface ApiResponse<T = any> {
//   code: number;
//   message: string;
//   data: T;
// }

// ==================== 请求参数类型 ====================

/**
 * 获取日统计数据请求参数
 * - data: 日期，格式为 YYYY-MM-DD
 * - type: 统计类型，可选值为 'users' | 'posts' | 'comments' | 'likes' | 'all'，默认为 'all'
 */
export interface GetStatsDayParams {
  date?: string;
  type?: "users" | "posts" | "comments" | "likes" | "all";
}

/**
 * 获取总计统计数据请求参数
 */
export interface GetStatsTotalParams {
  start_date?: string;
  end_date?: string;
  type?: "users" | "posts" | "comments" | "likes" | "all";
}

/**
 * 获取趋势统计数据请求参数
 */
export interface GetStatsTrendParams {
  start_date?: string;
  end_date?: string;
  type: "users" | "posts" | "comments" | "likes";
  interval?: "day" | "week" | "month";
}

// 范围统计请求参数
export interface GetStatsRangeParams {
  start_date?: string; // 2026-03-01
  end_date?: string; // 2026-03-31
  type?: "users" | "posts" | "comments" | "all";
}

// 每日统计数据项
export interface DailyStat {
  date: string; // 2026-03-21
  new_user: number;
  new_article: number;
  new_comment: number;
  new_board: number;
  new_tag: number;
  active_user: number;
}

// 范围统计响应
export type StatsRangeResponse = DailyStat[];
