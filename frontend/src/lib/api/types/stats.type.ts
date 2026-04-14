// stats.ts

/**
 * 系统基础统计信息
 */
export interface StatsInfo {
  /** 总用户数 */
  total_user: number;
  /** 总文章数 */
  total_article: number;
  /** 总评论数 */
  total_comment: number;
  /** 总板块数 */
  total_board: number;
  /** 总标签数 */
  total_tag: number;
}

/**
 * 今日统计信息
 */
export interface StatsTodayInfo {
  /** 今日新增用户 */
  new_user: number;
  /** 今日新增文章 */
  new_article: number;
  /** 今日新增评论 */
  new_comment: number;
  /** 今日新增板块 */
  new_board: number;
  /** 今日新增标签 */
  new_tag: number;
  /** 今日活跃用户数 */
  active_user: number;
}

/**
 * 违规统计信息
 */
export interface StatsIllegalInfo {
  /** 今日违规总数 */
  total: number;
  /** 今日违规用户数 */
  user_count: number;
  /** 今日违规文章数 */
  article_cnt: number;
  /** 今日违规评论数 */
  comment_cnt: number;
  /** 今日违规板块数 */
  board_cnt: number;
}

/**
 * 活跃用户详情
 */
export interface ActiveUserDetail {
  /** 用户ID */
  user_id: number;
  /** 用户名 */
  username: string;
  /** 头像 */
  avatar: string;
  /** 今日发文数 */
  article_count: number;
  /** 今日评论数 */
  comment_count: number;
  /** 最后活跃时间 */
  last_active_at: string; // ISO 8601 时间字符串
}

/**
 * 活跃用户信息
 */
export interface StatsActiveUserInfo {
  /** 今日活跃用户总数 */
  total: number;
  /** 今日活跃用户列表（最多N条） */
  list: ActiveUserDetail[];
}

/**
 * 热门文章项
 */
export interface HotArticleItem {
  /** 文章ID */
  id: number;
  /** 文章标题 */
  title: string;
  /** 板块ID */
  board_id: number;
  /** 板块名称 */
  board_name: string;
  /** 作者ID */
  author_id: number;
  /** 作者昵称 */
  author_name: string;
  /** 浏览量 */
  view_count: number;
  /** 评论数 */
  comment_count: number;
  /** 点赞数 */
  like_count: number;
  /** 综合热度分 */
  score: number;
}

/**
 * 热门板块项
 */
export interface HotBoardItem {
  /** 板块ID */
  id: number;
  /** 板块名称 */
  name: string;
  /** 板块图标 */
  icon: string;
  /** 今日发文数 */
  article_count: number;
  /** 今日评论数 */
  comment_count: number;
  /** 今日活跃用户数 */
  active_user: number;
  /** 综合热度分 */
  score: number;
}

/**
 * 违规用户项
 */
export interface StatsViolatorItem {
  /** 用户ID */
  user_id: number;
  /** 用户名 */
  username: string;
  /** 违规次数 */
  violation_cnt: number;
}

/**
 * 统计信息响应（聚合根）
 */
export interface StatsInfoResp {
  /** 基础统计信息 */
  base_info: StatsInfo;
  /** 今日统计信息 */
  today_info: StatsTodayInfo;
  /** 今日违规信息（可选） */
  illegal_info?: StatsIllegalInfo;
  /** 今日活跃用户信息（可选） */
  active_user_info?: StatsActiveUserInfo;
  /** 今日热门文章列表（可选） */
  hot_articles?: HotArticleItem[];
  /** 今日热门板块列表（可选） */
  hot_boards?: HotBoardItem[];
  /** 统计时间 */
  stat_time: string; // ISO 8601 时间字符串
}

/**
 * 日统计数据响应
 */
export interface StatsDayResponse {
  /** 日期 */
  day: string;
  /** 统计类型 */
  type: string;
  /** 数量 */
  count: number;
}

/**
 * 总计统计数据响应
 */
export interface StatsTotalResponse {
  /** 统计类型 */
  type: string;
  /** 数量 */
  count: number;
}

/**
 * 趋势数据
 */
export interface TrendData {
  /** 日期 */
  date: string;
  /** 数量 */
  count: number;
}

/**
 * 趋势统计响应
 */
export interface StatsTrendResponse {
  /** 开始日期 */
  start_date: string;
  /** 结束日期 */
  end_date: string;
  /** 统计粒度 (day/week/month) */
  interval: 'day' | 'week' | 'month';
  /** 统计类型 */
  type: 'users' | 'posts' | 'comments' | 'likes';
  /** 趋势数据 */
  trend: TrendData[];
}

/**
 * API 统一响应格式
 */
export interface ApiResponse<T = any> {
  /** 状态码，0 表示成功 */
  code: number;
  /** 响应消息 */
  message: string;
  /** 响应数据 */
  data: T;
}

// ==================== 请求参数类型 ====================

/**
 * 获取日统计数据请求参数
 */
export interface GetStatsDayParams {
  /** 日期 (格式: YYYY-MM-DD) */
  date?: string;
  /** 统计类型 */
  type?: 'users' | 'posts' | 'comments' | 'likes' | 'all';
}

/**
 * 获取总计统计数据请求参数
 */
export interface GetStatsTotalParams {
  /** 开始日期 (格式: YYYY-MM-DD) */
  start_date?: string;
  /** 结束日期 (格式: YYYY-MM-DD) */
  end_date?: string;
  /** 统计类型 */
  type?: 'users' | 'posts' | 'comments' | 'likes' | 'all';
}

/**
 * 获取趋势统计数据请求参数
 */
export interface GetStatsTrendParams {
  /** 开始日期 (格式: YYYY-MM-DD) */
  start_date?: string;
  /** 结束日期 (格式: YYYY-MM-DD) */
  end_date?: string;
  /** 统计类型 (必需) */
  type: 'users' | 'posts' | 'comments' | 'likes';
  /** 统计粒度 */
  interval?: 'day' | 'week' | 'month';
}