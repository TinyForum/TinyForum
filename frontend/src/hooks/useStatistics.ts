// hooks/useStatistics.ts
import { statsApi } from "@/lib/api/modules/stats";
import { GetStatsDayParams, GetStatsTotalParams, StatsTodayInfo, StatsInfoResp, StatsTrendResponse, GetStatsTrendParams } from "@/lib/api/types/stats.type";
import { useEffect, useState, useCallback } from "react";

interface UseStatisticsOptions {
  autoFetch?: boolean;     // 是否自动获取今日统计数据（默认 true）
  dayParams?: GetStatsDayParams;   // 今日统计参数，不传则 API 默认获取今日数据
  totalParams?: GetStatsTotalParams; // 总计统计参数（手动获取时使用）
}

interface UseStatisticsReturn {
  dayStats: StatsTodayInfo | null;
  dayLoading: boolean;
  dayError: Error | null;
  fetchDayStats: (params?: GetStatsDayParams) => Promise<void>;
  
  totalStats: StatsInfoResp | null;
  totalLoading: boolean;
  totalError: Error | null;
  fetchTotalStats: (params?: GetStatsTotalParams) => Promise<void>;
  
  trendStats: StatsTrendResponse | null;
  trendLoading: boolean;
  trendError: Error | null;
  fetchTrendStats: (params: GetStatsTrendParams) => Promise<void>;
  
  isLoading: boolean;
  refreshAllStatistics: () => Promise<void>; // 手动刷新今日+总计数据
}

/**
 * 获取统计数据
 * @param options 配置选项
 * - autoFetch 是否自动获取今日统计数据（默认 true），总计数据不会自动获取
 * - dayParams 获取今日统计数据时的参数，不传则默认获取今日数据
 * - totalParams 获取总统计数据时的参数（仅在手动调用 fetchTotalStats 或 refreshAllStatistics 时使用）
 * @returns 统计数据相关的状态和操作函数
 */
export const useStatistics = (options: UseStatisticsOptions = {}): UseStatisticsReturn => {
  const { autoFetch = true, dayParams, totalParams } = options;
  
  const [dayStats, setDayStats] = useState<StatsTodayInfo | null>(null);
  const [dayLoading, setDayLoading] = useState(false);
  const [dayError, setDayError] = useState<Error | null>(null);
  
  const [totalStats, setTotalStats] = useState<StatsInfoResp | null>(null);
  const [totalLoading, setTotalLoading] = useState(false);
  const [totalError, setTotalError] = useState<Error | null>(null);
  
  const [trendStats, setTrendStats] = useState<StatsTrendResponse | null>(null);
  const [trendLoading, setTrendLoading] = useState(false);
  const [trendError, setTrendError] = useState<Error | null>(null);
  
  const fetchDayStats = useCallback(async (params?: GetStatsDayParams) => {
    console.log("Fetching day stats with params:", params);
    setDayLoading(true);
    setDayError(null);
    try {
      const response = await statsApi.day(params);
      if (response.data?.data) {
        setDayStats(response.data.data);
      }
    } catch (error) {
      setDayError(error as Error);
      console.error("获取今日统计数据失败:", error);
    } finally {
      setDayLoading(false);
    }
  }, []);
  
  const fetchTotalStats = useCallback(async (params?: GetStatsTotalParams) => {
    setTotalLoading(true);
    setTotalError(null);
    try {
      const response = await statsApi.total(params);
      if (response.data?.data) {
        setTotalStats(response.data.data);
      }
    } catch (error) {
      setTotalError(error as Error);
      console.error("获取总计统计数据失败:", error);
    } finally {
      setTotalLoading(false);
    }
  }, []);
  
  const fetchTrendStats = useCallback(async (params: GetStatsTrendParams) => {
    setTrendLoading(true);
    setTrendError(null);
    try {
      const response = await statsApi.trend(params);
      if (response.data?.data) {
        setTrendStats(response.data.data);
      }
    } catch (error) {
      setTrendError(error as Error);
      console.error("获取趋势统计数据失败:", error);
    } finally {
      setTrendLoading(false);
    }
  }, []);
  
  const refreshAllStatistics = useCallback(async () => {
    // 手动刷新时同时获取今日和总计数据
    await Promise.all([
      fetchDayStats(dayParams),
      fetchTotalStats(totalParams),
    ]);
  }, [fetchDayStats, fetchTotalStats, dayParams, totalParams]);
  
  // 默认自动获取今日数据（不自动获取总计数据）
  useEffect(() => {
    if (autoFetch) {
      fetchDayStats(dayParams);
    }
  }, [autoFetch, fetchDayStats, dayParams]);
  
  return {
    dayStats,
    dayLoading,
    dayError,
    fetchDayStats,
    totalStats,
    totalLoading,
    totalError,
    fetchTotalStats,
    trendStats,
    trendLoading,
    trendError,
    fetchTrendStats,
    isLoading: dayLoading || totalLoading || trendLoading,
    refreshAllStatistics,
  };
};