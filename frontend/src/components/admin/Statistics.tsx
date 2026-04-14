import { useEffect, useRef } from "react";
import * as echarts from "echarts";
import { useStatsData } from "@/hooks/admin/useStatsData";
import { Activity, Award, Database, FileText, Users } from "lucide-react";

export function Statistics({ t }: { t: (key: string) => string }) {
  const { stats, charts, exportData, isLoading } = useStatsData(true);

  const userChartRef = useRef<HTMLDivElement>(null);
  const postChartRef = useRef<HTMLDivElement>(null);

  // 渲染用户增长图表
  useEffect(() => {
    if (!userChartRef.current || !charts?.userGrowthTrend?.length) return;

    const chart = echarts.init(userChartRef.current);
    chart.setOption({
      tooltip: { trigger: "axis" },
      xAxis: {
        type: "category",
        data: charts.userGrowthTrend.map((item) => item.date.slice(5)),
      },
      yAxis: { type: "value", name: t("user_count") },
      series: [
        {
          data: charts.userGrowthTrend.map((item) => item.value),
          type: "line",
          smooth: true,
          areaStyle: { opacity: 0.3 },
          lineStyle: { color: "#6366f1", width: 2 },
          itemStyle: { color: "#6366f1" },
        },
      ],
      grid: { top: 30, left: 50, right: 20, bottom: 20, containLabel: true },
    });

    return () => chart.dispose();
  }, [charts?.userGrowthTrend, t]);

  // 渲染帖子增长图表
  useEffect(() => {
    if (!postChartRef.current || !charts?.postTrend?.length) return;

    const chart = echarts.init(postChartRef.current);
    chart.setOption({
      tooltip: { trigger: "axis" },
      xAxis: {
        type: "category",
        data: charts.postTrend.map((item) => item.date.slice(5)),
      },
      yAxis: { type: "value", name: t("post_count") },
      series: [
        {
          data: charts.postTrend.map((item) => item.value),
          type: "line",
          smooth: true,
          areaStyle: { opacity: 0.3 },
          lineStyle: { color: "#06b6d4", width: 2 },
          itemStyle: { color: "#06b6d4" },
        },
      ],
      grid: { top: 30, left: 50, right: 20, bottom: 20, containLabel: true },
    });

    return () => chart.dispose();
  }, [charts?.postTrend, t]);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 核心指标卡片 */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-primary">
            <Users className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("total_users")}</div>
          <div className="stat-value text-2xl md:text-3xl">
            {stats?.totalUsers?.toLocaleString() || 0}
          </div>
          <div className="stat-desc text-success">
            ↑ {stats?.userGrowth || 0}%
          </div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-secondary">
            <FileText className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("total_posts")}</div>
          <div className="stat-value text-2xl md:text-3xl">
            {stats?.totalPosts?.toLocaleString() || 0}
          </div>
          <div className="stat-desc text-success">
            ↑ {stats?.postGrowth || 0}%
          </div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-accent">
            <Activity className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("today_active")}</div>
          <div className="stat-value text-2xl md:text-3xl">
            {stats?.todayActive?.toLocaleString() || 0}
          </div>
          <div className="stat-desc">
            {t("online_now")}:{" "}
            <span className="text-success">{stats?.onlineNow || 0}</span>
          </div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-warning">
            <Award className="w-6 h-6" />
          </div>
          <div className="stat-title">{t("total_points")}</div>
          <div className="stat-value text-2xl md:text-3xl">
            {stats?.totalPoints?.toLocaleString() || 0}
          </div>
        </div>
      </div>

      {/* 图表区域 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body">
            <h3 className="font-semibold text-lg">{t("user_growth_trend")}</h3>
            <div ref={userChartRef} className="h-64 w-full" />
          </div>
        </div>

        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body">
            <h3 className="font-semibold text-lg">{t("post_trend")}</h3>
            <div ref={postChartRef} className="h-64 w-full" />
          </div>
        </div>
      </div>

      {/* 操作栏 */}
      <div className="flex justify-end gap-2">
        <button
          onClick={() => window.location.reload()}
          className="btn btn-ghost btn-sm"
        >
          🔄 {t("refresh")}
        </button>
        <button
          onClick={() => exportData("csv")}
          className="btn btn-primary btn-sm gap-2"
        >
          <Database className="w-4 h-4" />
          {t("export_data")}
        </button>
      </div>
    </div>
  );
}
