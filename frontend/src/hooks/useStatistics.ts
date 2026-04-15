// hooks/useStatistics.ts
import { statsApi } from "@/lib/api/modules/stats";
import { GetStatsDayParams, GetStatsTotalParams, StatsTodayInfo, StatsInfoResp, StatsTrendResponse, GetStatsTrendParams } from "@/lib/api/types/stats.type";
import { useEffect, useState, useCallback } from "react";


interface UseStatisticsOptions {
  autoFetch?: boolean;
  dayParams?: GetStatsDayParams;
  totalParams?: GetStatsTotalParams;
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
  refreshAllStatistics: () => Promise<void>;
}

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
    await Promise.all([
      fetchDayStats(dayParams),
      fetchTotalStats(totalParams),
    ]);
  }, [fetchDayStats, fetchTotalStats, dayParams, totalParams]);
  
  useEffect(() => {
    if (autoFetch) {
      refreshAllStatistics();
    }
  }, [autoFetch, refreshAllStatistics]);
  
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