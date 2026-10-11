// hooks/admin/useStatsData.ts

import { statsApi } from "@/shared/api/modules/stats";
// import { StatsInfoResp } from "@/shared/api/types/stats.type";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useState, useCallback } from "react";
import toast from "react-hot-toast";
import { useTranslations } from "next-intl";
import { statsKeys } from "@/features/statistic/hooks/useStatsKeys";

// 依据趋势序列最后两点计算环比增长率（%），数据不足或基数为 0 时返回 0
function calcGrowthFromTrend(
  trend?: Array<{ date: string; count: number }>,
): number {
  if (!trend || trend.length < 2) return 0;
  const previous = trend[trend.length - 2]?.count ?? 0;
  const current = trend[trend.length - 1]?.count ?? 0;
  if (previous === 0) return 0;
  return parseFloat((((current - previous) / previous) * 100).toFixed(1));
}

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
  const t = useTranslations("Admin");

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

  // 1. 获取总计统计数据（包含基础信息和今日统计）
  const {
    data: totalData,
    isLoading: totalLoading,
    // refetchTotal, // 移除未使用的变量
  } = useQuery({
    queryKey: statsKeys.totals(dateRange.start, dateRange.end),
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
    queryKey: statsKeys.trend("users", dateRange.start, dateRange.end),
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
    queryKey: statsKeys.trend("posts", dateRange.start, dateRange.end),
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

  // 处理统计数据 - 修复可能为 undefined 的问题
  const stats = (() => {
    if (!totalData?.data?.data) return null;

    const todayInfo = totalData.data.data.today_info;
    const baseInfo = totalData.data.data.base_info;

    const userTrendSeries = userTrend?.data?.data?.trend || [];
    const postTrendSeries = postTrend?.data?.data?.trend || [];

    return {
      totalUsers: baseInfo?.total_user || 0,
      userGrowth: calcGrowthFromTrend(userTrendSeries),
      totalPosts: baseInfo?.total_article || 0,
      postGrowth: calcGrowthFromTrend(postTrendSeries),
      todayActive: todayInfo?.active_user || 0,
      onlineNow: 0,
      totalPoints: 0,
    };
  })();

  // 处理图表数据 - 修复可能为 undefined 的问题
  const charts = (() => {
    return {
      userGrowthTrend:
        userTrend?.data?.data?.trend?.map(
          (item: { date: string; count: number }) => ({
            date: item.date,
            value: item.count,
          }),
        ) || [],
      postTrend:
        postTrend?.data?.data?.trend?.map(
          (item: { date: string; count: number }) => ({
            date: item.date,
            value: item.count,
          }),
        ) || [],
    };
  })();

  // 导出数据：从真实趋势序列生成 CSV
  const exportData = useCallback(
    async (format: "csv" | "excel" = "csv") => {
      try {
        const userSeries = userTrend?.data?.data?.trend || [];
        const postSeries = postTrend?.data?.data?.trend || [];
        const userMap = new Map(
          userSeries.map((i: { date: string; count: number }) => [
            i.date,
            i.count,
          ]),
        );
        const postMap = new Map(
          postSeries.map((i: { date: string; count: number }) => [
            i.date,
            i.count,
          ]),
        );
        const dates = Array.from(
          new Set([...userMap.keys(), ...postMap.keys()]),
        ).sort();
        const rows: string[][] = [
          ["date", "new_users", "new_posts"],
          ...dates.map((date) => [
            date,
            String(userMap.get(date) ?? 0),
            String(postMap.get(date) ?? 0),
          ]),
        ];
        const csv = rows
          .map((row) =>
            row.map((cell) => `"${cell.replace(/"/g, '""')}"`).join(","),
          )
          .join("\n");
        const blob = new Blob(["\uFEFF" + csv], {
          type: "text/csv;charset=utf-8;",
        });
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.download = `statistics_${dateRange.start}_${dateRange.end}.${format === "excel" ? "xls" : "csv"}`;
        link.click();
        URL.revokeObjectURL(url);
      } catch {
        toast.error(t("export_failed"));
      }
    },
    [dateRange, userTrend, postTrend, t],
  );

  // 刷新所有数据
  const refreshAll = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: statsKeys.all });
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

    // 原始数据（用于调试）
    rawData: {
      totalData: totalData?.data,
      userTrend: userTrend?.data,
      postTrend: postTrend?.data,
    },
  };
}
