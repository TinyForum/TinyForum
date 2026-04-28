import apiClient from "../client";
import { ApiResponse } from "../types";
import {
  // StatsDayResponse,
  StatsTodayInfo,
  StatsInfoResp,
  StatsTrendResponse,
  GetStatsDayParams,
  GetStatsTotalParams,
  GetStatsTrendParams,
  GetStatsRangeParams,
  StatsRangeResponse,
} from "@/shared/api/types/stats.type";

export const statsApi = {
  // ── 统计数据管理 ──────────────────────────────────────────────────────────────

  /**
   * 获取日统计数据
   * @param params - 查询参数
   * @returns 今日统计数据
   */
  day: (params?: GetStatsDayParams) =>
    apiClient.get<ApiResponse<StatsTodayInfo>>("/statistics/day", { params }),

  /**
   * 获取总计统计数据
   * @param params - 查询参数
   * @returns 总计统计数据（包含基础信息、今日信息、热门内容等）
   */
  total: (params?: GetStatsTotalParams) =>
    apiClient.get<ApiResponse<StatsInfoResp>>("/statistics/total", { params }),

  /**
   * 获取趋势统计数据
   * @param params - 查询参数
   * @returns 趋势数据（按天/周/月统计）
   */
  trend: (params: GetStatsTrendParams) =>
    apiClient.get<ApiResponse<StatsTrendResponse>>("/statistics/trend", {
      params,
    }),
  /**
   * 获取范围统计数据
   */
  range: (params?: GetStatsRangeParams) =>
    apiClient.get<ApiResponse<StatsRangeResponse>>("/statistics/range", {
      params,
    }),
};
