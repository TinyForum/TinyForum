// hooks/useStatistics.ts
import { statsApi } from "@/shared/api/modules/stats";
import {
  GetStatsDayParams,
  GetStatsTotalParams,
  StatsTodayInfo,
  StatsInfoResp,
  GetStatsRangeParams,
  StatsRangeResponse,
} from "@/shared/api/types/stats.model";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { statsKeys } from "./hooks/useStatsKeys";

interface UseStatisticsOptions {
  autoFetch?: boolean; // 是否自动获取今日/总计统计数据（默认 true）
  dayParams?: GetStatsDayParams; // 今日统计参数
  totalParams?: GetStatsTotalParams; // 总计统计参数
  rangeParams?: GetStatsRangeParams; // 范围统计参数（传入即自动获取）
}

interface UseStatisticsReturn {
  // 今日统计
  dayStats: StatsTodayInfo | null;
  dayLoading: boolean;
  // 总计统计
  totalStats: StatsInfoResp | null;
  totalLoading: boolean;
  // 范围统计
  rangeStats: StatsRangeResponse | null;
  rangeLoading: boolean;
  isLoading: boolean;
  refreshAllStatistics: () => Promise<void>; // 刷新今日+总计+范围
}

const unwrap = <T>(data: T | null | undefined): T | null => data ?? null;

export const useStatistics = (
  options: UseStatisticsOptions = {},
): UseStatisticsReturn => {
  const { autoFetch = true, dayParams, totalParams, rangeParams } = options;
  const queryClient = useQueryClient();

  // 今日统计
  const dayQuery = useQuery<StatsTodayInfo>({
    queryKey: statsKeys.day(dayParams),
    queryFn: async () => {
      const res = await statsApi.day(dayParams);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取今日统计数据失败");
      }
      return res.data.data!;
    },
    enabled: autoFetch,
  });

  // 总计统计
  const totalQuery = useQuery<StatsInfoResp>({
    queryKey: statsKeys.total(totalParams),
    queryFn: async () => {
      const res = await statsApi.total(totalParams);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取总计统计数据失败");
      }
      return res.data.data!;
    },
    enabled: autoFetch,
  });

  // 范围统计（传入 rangeParams 即自动获取）
  const rangeQuery = useQuery<StatsRangeResponse>({
    queryKey: statsKeys.range(rangeParams),
    queryFn: async () => {
      const res = await statsApi.range(rangeParams);
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取范围统计数据失败");
      }
      return res.data.data!;
    },
    enabled: !!rangeParams,
  });

  // 刷新今日+总计+范围
  const refreshAllStatistics = async () => {
    await Promise.all([
      queryClient.refetchQueries({ queryKey: statsKeys.day(dayParams) }),
      queryClient.refetchQueries({ queryKey: statsKeys.total(totalParams) }),
      rangeParams &&
        queryClient.refetchQueries({ queryKey: statsKeys.range(rangeParams) }),
    ]);
  };

  return {
    dayStats: unwrap(dayQuery.data),
    dayLoading: dayQuery.isLoading,
    totalStats: unwrap(totalQuery.data),
    totalLoading: totalQuery.isLoading,
    rangeStats: unwrap(rangeQuery.data),
    rangeLoading: rangeQuery.isLoading,
    isLoading: dayQuery.isLoading || totalQuery.isLoading || rangeQuery.isLoading,
    refreshAllStatistics,
  };
};
