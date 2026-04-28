import { useEffect, useRef } from "react";
import * as echarts from "echarts";
import { useStatsData } from "@/hooks/admin/useStatsData";
import { useStatistics } from "@/hooks/useStatistics";
import {
  Activity,
  Award,
  Database,
  FileText,
  Users,
  MessageCircle,
  TrendingUp,
} from "lucide-react";
import { CommunityStats } from "../home/CommunityStats";

export function Dashboard({ t }: { t: (key: string) => string }) {
  const { stats,  exportData, isLoading } = useStatsData(true);

  // 获取范围统计数据（最近30天，type=all）
  const { rangeStats,  fetchRangeStats } = useStatistics({
    autoFetch: false,
  });

  // 图表 refs（原有 + 新增）
  const newUserTrendRef = useRef<HTMLDivElement>(null);
  const newPostTrendRef = useRef<HTMLDivElement>(null);
  const newCommentTrendRef = useRef<HTMLDivElement>(null);
  const allTrendRef = useRef<HTMLDivElement>(null);

  // 获取范围数据（默认最近30天）
  useEffect(() => {
    fetchRangeStats({ type: "all" });
  }, [fetchRangeStats]);

  // ========== 新增：新增用户趋势图 ==========
  useEffect(() => {
    if (!newUserTrendRef.current || !rangeStats?.length) return;
    const chart = echarts.init(newUserTrendRef.current);
    chart.setOption({
      tooltip: { trigger: "axis" },
      xAxis: { type: "category", data: rangeStats.map((d) => d.date.slice(5)) },
      yAxis: { type: "value", name: t("user_count") },
      series: [
        {
          data: rangeStats.map((d) => d.new_user),
          type: "line",
          smooth: true,
          areaStyle: { opacity: 0.2 },
          lineStyle: { color: "#6366f1", width: 2 },
        },
      ],
      grid: { top: 30, left: 50, right: 20, bottom: 20, containLabel: true },
    });
    return () => chart.dispose();
  }, [rangeStats, t]);

  // ========== 新增：新增文章趋势图 ==========
  useEffect(() => {
    if (!newPostTrendRef.current || !rangeStats?.length) return;
    const chart = echarts.init(newPostTrendRef.current);
    chart.setOption({
      tooltip: { trigger: "axis" },
      xAxis: { type: "category", data: rangeStats.map((d) => d.date.slice(5)) },
      yAxis: { type: "value", name: t("post_count") },
      series: [
        {
          data: rangeStats.map((d) => d.new_article),
          type: "line",
          smooth: true,
          areaStyle: { opacity: 0.2 },
          lineStyle: { color: "#06b6d4", width: 2 },
        },
      ],
      grid: { top: 30, left: 50, right: 20, bottom: 20, containLabel: true },
    });
    return () => chart.dispose();
  }, [rangeStats, t]);

  // ========== 新增：新增评论趋势图 ==========
  useEffect(() => {
    if (!newCommentTrendRef.current || !rangeStats?.length) return;
    const chart = echarts.init(newCommentTrendRef.current);
    chart.setOption({
      tooltip: { trigger: "axis" },
      xAxis: { type: "category", data: rangeStats.map((d) => d.date.slice(5)) },
      yAxis: { type: "value", name: t("comment_count") },
      series: [
        {
          data: rangeStats.map((d) => d.new_comment),
          type: "line",
          smooth: true,
          areaStyle: { opacity: 0.2 },
          lineStyle: { color: "#ec489a", width: 2 },
        },
      ],
      grid: { top: 30, left: 50, right: 20, bottom: 20, containLabel: true },
    });
    return () => chart.dispose();
  }, [rangeStats, t]);

  // ========== 新增：综合趋势（文章+评论） ==========
  useEffect(() => {
    if (!allTrendRef.current || !rangeStats?.length) return;
    const chart = echarts.init(allTrendRef.current);
    const totalNew = rangeStats.map((d) => d.new_article + d.new_comment);
    chart.setOption({
      tooltip: { trigger: "axis" },
      xAxis: { type: "category", data: rangeStats.map((d) => d.date.slice(5)) },
      yAxis: { type: "value", name: t("total_new_content") },
      series: [
        {
          data: totalNew,
          type: "line",
          smooth: true,
          areaStyle: { opacity: 0.2 },
          lineStyle: { color: "#f59e0b", width: 2 },
        },
      ],
      grid: { top: 30, left: 50, right: 20, bottom: 20, containLabel: true },
    });
    return () => chart.dispose();
  }, [rangeStats, t]);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 核心指标卡片（原有，保持不变） */}
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

      {/* 新增四个趋势图 */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1">
              <Users className="w-4 h-4" /> {t("new_users_trend")}
            </h4>
            <div ref={newUserTrendRef} className="h-48 w-full" />
          </div>
        </div>
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1">
              <FileText className="w-4 h-4" /> {t("new_posts_trend")}
            </h4>
            <div ref={newPostTrendRef} className="h-48 w-full" />
          </div>
        </div>
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1">
              <MessageCircle className="w-4 h-4" /> {t("new_comments_trend")}
            </h4>
            <div ref={newCommentTrendRef} className="h-48 w-full" />
          </div>
        </div>
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1">
              <TrendingUp className="w-4 h-4" /> {t("total_new_content")}
            </h4>
            <div ref={allTrendRef} className="h-48 w-full" />
          </div>
        </div>
      </div>

      {/* 社区统计组件（原有） */}
      <CommunityStats />

      {/* 操作栏（原有） */}
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
