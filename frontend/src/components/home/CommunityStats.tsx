// components/home/CommunityStats.tsx
import { useStatistics } from "@/hooks/useStatistics";
import { Sparkles, TrendingUp } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useRef, } from "react";
import * as echarts from "echarts";
import { SkeletonItem } from "@/shared/ui/SkeletonItem";
import { StatItem } from "@/shared/ui/StatItem";

interface CommunityStatsProps {
  className?: string;
}

// 辅助函数移到组件外部
const getToday = (): string => {
  const d = new Date();
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
};

const getLast7DaysStart = (): string => {
  const d = new Date();
  d.setDate(d.getDate() - 6);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, "0")}-${String(d.getDate()).padStart(2, "0")}`;
};

export const CommunityStats = ({ className = "" }: CommunityStatsProps) => {
  const t = useTranslations("Sidebar");
  const {
    dayStats,
    totalStats,
    rangeStats,
    rangeLoading,
    fetchRangeStats,
    isLoading,
  } = useStatistics({
    autoFetch: true,
  });

  const chartRef = useRef<HTMLDivElement>(null);

  // 获取最近7天的范围统计（仅帖子数量）
  useEffect(() => {
    fetchRangeStats({
      start_date: getLast7DaysStart(),
      end_date: getToday(),
      type: "posts",
    });
  }, [fetchRangeStats]);

  // 渲染迷你趋势图
  useEffect(() => {
    if (!chartRef.current || !rangeStats || rangeStats.length === 0) return;

    const chart = echarts.init(chartRef.current);
    
    chart.setOption({
      tooltip: { trigger: "axis", axisPointer: { type: "shadow" } },
      grid: { top: 20, left: 35, right: 5, bottom: 5, containLabel: true },
      xAxis: {
        type: "category",
        data: rangeStats.map((item: { date: string }) => item.date.slice(5)),
        axisLabel: { rotate: 30, fontSize: 10 },
      },
      yAxis: {
        type: "value",
        name: t("post_count"),
        nameTextStyle: { fontSize: 10 },
        splitLine: { lineStyle: { type: "dashed" } },
      },
      series: [
        {
          data: rangeStats.map((item: { new_article: number }) => item.new_article),
          type: "line",
          smooth: true,
          lineStyle: { color: "#06b6d4", width: 2 },
          areaStyle: { opacity: 0.2, color: "#06b6d4" },
          symbol: "circle",
          symbolSize: 6,
          itemStyle: { color: "#0891b2" },
        },
      ],
    });

    const handleResize = () => chart.resize();
    window.addEventListener("resize", handleResize);

    return () => {
      window.removeEventListener("resize", handleResize);
      chart.dispose();
    };
  }, [rangeStats, t]);

  const totalUsers = totalStats?.base_info?.total_user || 0;
  const totalPosts = totalStats?.base_info?.total_article || 0;
  const todayPosts = dayStats?.new_article || 0;
  const yesterdayActive = dayStats?.active_user || 0;

  return (
    <div className={`rounded-lg border bg-card shadow-sm ${className}`}>
      <div className="p-3 border-b">
        <h3 className="font-semibold flex items-center gap-2">
          <Sparkles className="w-4 h-4" />
          {t("community_stats")}
        </h3>
      </div>
      <div className="p-3 space-y-3 text-sm">
        {isLoading ? (
          <>
            <SkeletonItem />
            <SkeletonItem />
            <SkeletonItem />
            <SkeletonItem />
          </>
        ) : (
          <>
            <StatItem label={t("today_posts")} value={todayPosts} />
            <StatItem label={t("yesterday_active")} value={yesterdayActive} />
            <StatItem label={t("total_users")} value={totalUsers} />
            <StatItem label={t("total_posts")} value={totalPosts} />
          </>
        )}

        {/* 趋势图区域（仅当有数据时显示） */}
        {!rangeLoading && rangeStats && rangeStats.length > 0 && (
          <div className="pt-2 border-t border-border">
            <div className="flex items-center gap-1 text-xs text-muted-foreground mb-2">
              <TrendingUp className="w-3 h-3" />
              <span>{t("last_7_days_posts_trend")}</span>
            </div>
            <div ref={chartRef} className="h-24 w-full" />
          </div>
        )}
      </div>
    </div>
  );
};