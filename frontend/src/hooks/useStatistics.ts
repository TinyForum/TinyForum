// hooks/useStatistics.ts
import { statsApi } from "@/lib/api/modules/stats";
import {
  GetStatsDayParams,
  GetStatsTotalParams,
  StatsTodayInfo,
  StatsInfoResp,
  StatsTrendResponse,
  GetStatsTrendParams,
  GetStatsRangeParams,    // 新增导入
  StatsRangeResponse,     // 新增导入
} from "@/lib/api/types/stats.type";
import { useEffect, useState, useCallback } from "react";

interface UseStatisticsOptions {
  autoFetch?: boolean;               // 是否自动获取今日统计数据（默认 true）
  dayParams?: GetStatsDayParams;     // 今日统计参数
  totalParams?: GetStatsTotalParams; // 总计统计参数
  // 注意：range 参数需要手动传入 fetchRangeStats 调用，不自动获取
}

interface UseStatisticsReturn {
  // 原有
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

  // 新增范围统计
  rangeStats: StatsRangeResponse | null;
  rangeLoading: boolean;
  rangeError: Error | null;
  fetchRangeStats: (params?: GetStatsRangeParams) => Promise<void>;

  isLoading: boolean;
  refreshAllStatistics: () => Promise<void>; // 刷新今日+总计（不包含范围）
}

export const useStatistics = (options: UseStatisticsOptions = {}): UseStatisticsReturn => {
  const { autoFetch = true, dayParams, totalParams } = options;

  // 原有状态
  const [dayStats, setDayStats] = useState<StatsTodayInfo | null>(null);
  const [dayLoading, setDayLoading] = useState(false);
  const [dayError, setDayError] = useState<Error | null>(null);

  const [totalStats, setTotalStats] = useState<StatsInfoResp | null>(null);
  const [totalLoading, setTotalLoading] = useState(false);
  const [totalError, setTotalError] = useState<Error | null>(null);

  const [trendStats, setTrendStats] = useState<StatsTrendResponse | null>(null);
  const [trendLoading, setTrendLoading] = useState(false);
  const [trendError, setTrendError] = useState<Error | null>(null);

  // 新增范围统计状态
  const [rangeStats, setRangeStats] = useState<StatsRangeResponse | null>(null);
  const [rangeLoading, setRangeLoading] = useState(false);
  const [rangeError, setRangeError] = useState<Error | null>(null);

  // 原有方法（省略具体实现，保持原样）
  const fetchDayStats = useCallback(async (params?: GetStatsDayParams) => {
    setDayLoading(true);
    setDayError(null);
    try {
      const response = await statsApi.day(params);
      if (response.data?.data) setDayStats(response.data.data);
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
      if (response.data?.data) setTotalStats(response.data.data);
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
      if (response.data?.data) setTrendStats(response.data.data);
    } catch (error) {
      setTrendError(error as Error);
      console.error("获取趋势统计数据失败:", error);
    } finally {
      setTrendLoading(false);
    }
  }, []);

  // 新增 fetchRangeStats 方法
  const fetchRangeStats = useCallback(async (params?: GetStatsRangeParams) => {
    setRangeLoading(true);
    setRangeError(null);
    try {
      const response = await statsApi.range(params);
      if (response.data?.data) {
        setRangeStats(response.data.data);
      }
    } catch (error) {
      setRangeError(error as Error);
      console.error("获取范围统计数据失败:", error);
    } finally {
      setRangeLoading(false);
    }
  }, []);

  // 刷新今日+总计（不包含范围，因为范围需要单独参数）
  const refreshAllStatistics = useCallback(async () => {
    await Promise.all([fetchDayStats(dayParams), fetchTotalStats(totalParams)]);
  }, [fetchDayStats, fetchTotalStats, dayParams, totalParams]);

  // 自动获取今日数据（原有逻辑）
  useEffect(() => {
    if (autoFetch) {
      fetchDayStats(dayParams);
    }
  }, [autoFetch, fetchDayStats, dayParams]);

  return {
    // 原有
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
    // 新增
    rangeStats,
    rangeLoading,
    rangeError,
    fetchRangeStats,
    isLoading: dayLoading || totalLoading || trendLoading || rangeLoading,
    refreshAllStatistics,
  };
};