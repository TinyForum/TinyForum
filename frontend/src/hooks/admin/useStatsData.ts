// hooks/admin/useStatsData.ts

import { statsApi } from "@/lib/api/modules/stats";
// import { StatsInfoResp } from "@/lib/api/types/stats.type";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useState, useEffect, useCallback } from "react";

// 注释掉未使用的接口
// interface UseStatsDataOptions {
//   enabled?: boolean;
//   defaultDays?: number;
// }

// interface StatsData {
//   stats: {
//     totalUsers: number;
//     userGrowth: number;
//     totalPosts: number;
//     postGrowth: number;
//     todayActive: number;
//     onlineNow: number;
//     totalPoints: number;
//   } | null;
//   charts: {
//     userGrowthTrend: Array<{ date: string; value: number }>;
//     postTrend: Array<{ date: string; value: number }>;
//   } | null;
//   hotArticles: StatsInfoResp["hot_articles"];
//   hotBoards: StatsInfoResp["hot_boards"];
//   activeUsers: StatsInfoResp["active_user_info"];
// }

export function useStatsData(enabled: boolean = true) {
  const queryClient = useQueryClient();

  // 日期范围状态
  const [dateRange, setDateRange] = useState<{ start: string; end: string }>(
    () => {
      const end = new Date();
      const start = new Date();
      start.setDate(start.getDate() - 30);
      return {
        start: start.toISOString().split("T")[0],
        end: end.toISOString().split("T")[0],
      };
    },
  );

  // 在线人数状态
  const [onlineNow, setOnlineNow] = useState(0);

  // 1. 获取总计统计数据（包含基础信息和今日统计）
  const {
    data: totalData,
    isLoading: totalLoading,
    // refetchTotal, // 移除未使用的变量
  } = useQuery({
    queryKey: ["stats", "total", dateRange.start, dateRange.end],
    queryFn: () =>
      statsApi.total({
        start_date: dateRange.start,
        end_date: dateRange.end,
        type: "all",
      }),
    enabled,
    staleTime: 5 * 60 * 1000,
  });

  // 2. 获取用户增长趋势
  const { data: userTrend, isLoading: userTrendLoading } = useQuery({
    queryKey: ["stats", "trend", "users", dateRange.start, dateRange.end],
    queryFn: () =>
      statsApi.trend({
        start_date: dateRange.start,
        end_date: dateRange.end,
        type: "users",
        interval: "day",
      }),
    enabled,
    staleTime: 10 * 60 * 1000,
  });

  // 3. 获取帖子增长趋势
  const { data: postTrend, isLoading: postTrendLoading } = useQuery({
    queryKey: ["stats", "trend", "posts", dateRange.start, dateRange.end],
    queryFn: () =>
      statsApi.trend({
        start_date: dateRange.start,
        end_date: dateRange.end,
        type: "posts",
        interval: "day",
      }),
    enabled,
    staleTime: 10 * 60 * 1000,
  });

  // 模拟在线人数（实际应接入 WebSocket 或轮询接口）
  useEffect(() => {
    if (!enabled) return;

    const fetchOnlineCount = async () => {
      setOnlineNow(0);
    };

    fetchOnlineCount();
    const interval = setInterval(fetchOnlineCount, 30000);

    return () => clearInterval(interval);
  }, [enabled]);

  // 计算增长率（对比昨日）
  const calculateGrowth = useCallback(
    (current: number, previous: number): number => {
      if (previous === 0) return 0;
      return parseFloat((((current - previous) / previous) * 100).toFixed(1));
    },
    [],
  );

  // 处理统计数据 - 修复可能为 undefined 的问题
  const stats = (() => {
    if (!totalData?.data?.data) return null;

    const todayInfo = totalData.data.data.today_info;
    const baseInfo = totalData.data.data.base_info;

    // 为了计算增长率，这里使用模拟数据
    const previousNewUser = Math.floor((todayInfo?.new_user || 0) / 1.2);
    const previousNewArticle = Math.floor((todayInfo?.new_article || 0) / 1.15);

    return {
      totalUsers: baseInfo?.total_user || 0,
      userGrowth: calculateGrowth(todayInfo?.new_user || 0, previousNewUser),
      totalPosts: baseInfo?.total_article || 0,
      postGrowth: calculateGrowth(
        todayInfo?.new_article || 0,
        previousNewArticle,
      ),
      todayActive: todayInfo?.active_user || 0,
      onlineNow,
      totalPoints: 0,
    };
  })();

  // 处理图表数据 - 修复可能为 undefined 的问题
  const charts = (() => {
    return {
      userGrowthTrend:
        userTrend?.data?.data?.trend?.map((item: { date: string; count: number }) => ({
          date: item.date,
          value: item.count,
        })) || [],
      postTrend:
        postTrend?.data?.data?.trend?.map((item: { date: string; count: number }) => ({
          date: item.date,
          value: item.count,
        })) || [],
    };
  })();

  // 导出数据
  const exportData = useCallback(
    async (format: "csv" | "excel" = "csv") => {
      try {
        console.log("导出数据:", { dateRange, format });
        const blob = new Blob(["模拟数据"], { type: "text/csv" });
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.download = `statistics_${dateRange.start}_${dateRange.end}.${format}`;
        link.click();
        URL.revokeObjectURL(url);
      } catch (error) {
        console.error("导出失败:", error);
      }
    },
    [dateRange],
  );

  // 刷新所有数据
  const refreshAll = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: ["stats"] });
  }, [queryClient]);

  // 更新日期范围
  const updateDateRange = useCallback((start: string, end: string) => {
    setDateRange({ start, end });
  }, []);

  // 预设时间范围
  const setPresetRange = useCallback(
    (range: "today" | "week" | "month" | "year") => {
      const end = new Date();
      const start = new Date();

      switch (range) {
        case "today":
          updateDateRange(
            end.toISOString().split("T")[0],
            end.toISOString().split("T")[0],
          );
          break;
        case "week":
          start.setDate(start.getDate() - 7);
          updateDateRange(
            start.toISOString().split("T")[0],
            end.toISOString().split("T")[0],
          );
          break;
        case "month":
          start.setMonth(start.getMonth() - 1);
          updateDateRange(
            start.toISOString().split("T")[0],
            end.toISOString().split("T")[0],
          );
          break;
        case "year":
          start.setFullYear(start.getFullYear() - 1);
          updateDateRange(
            start.toISOString().split("T")[0],
            end.toISOString().split("T")[0],
          );
          break;
      }
    },
    [updateDateRange],
  );

  const isLoading = totalLoading || userTrendLoading || postTrendLoading;

  // 安全获取数据的辅助函数
  const getTotalData = () => totalData?.data?.data;
  const totalDataValue = getTotalData();

  return {
    // 核心数据
    stats,
    charts,

    // 扩展数据 - 使用可选链安全访问
    hotArticles: totalDataValue?.hot_articles || [],
    hotBoards: totalDataValue?.hot_boards || [],
    activeUsers: totalDataValue?.active_user_info,
    illegalInfo: totalDataValue?.illegal_info,
    baseInfo: totalDataValue?.base_info,
    todayInfo: totalDataValue?.today_info,

    // 状态
    isLoading,
    isError: !totalData && !isLoading,

    // 操作方法
    exportData,
    refreshAll,
    updateDateRange,
    setPresetRange,

    // 当前参数
    dateRange,
    onlineNow,

    // 原始数据（用于调试）
    rawData: {
      totalData: totalData?.data,
      userTrend: userTrend?.data,
      postTrend: postTrend?.data,
    },
  };
}