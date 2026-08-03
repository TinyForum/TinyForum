"use client";

import { useEffect, useRef } from "react";
import * as echarts from "echarts";
import {
  useRecOverview,
  useRecBehaviorStats,
  useRecUserAnalysis,
  useRecContentPerformance,
  useRecRiskAnalysis,
} from "../hooks/useRecommendationStats";
import {
  Activity,
  Eye,
  MousePointerClick,
  Target,
  TrendingUp,
  Users,
  Shield,
  AlertTriangle,
  Ban,
  Flag,
} from "lucide-react";

export function RecommendationAnalytics() {
  const { data: overview, isLoading } = useRecOverview(true);
  const { data: behaviorStats, isLoading: bhLoading } = useRecBehaviorStats(7, true);
  const { data: userAnalysis, isLoading: uaLoading } = useRecUserAnalysis(7, true);
  const { data: contentPerf, isLoading: cpLoading } = useRecContentPerformance(true);
  const { data: riskAnalysis, isLoading: riskLoading } = useRecRiskAnalysis(true);

  const behaviorDistRef = useRef<HTMLDivElement>(null);
  const dailyTrendRef = useRef<HTMLDivElement>(null);
  const qualityDistRef = useRef<HTMLDivElement>(null);
  const boardDistRef = useRef<HTMLDivElement>(null);
  const tagCloudRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!behaviorDistRef.current || !behaviorStats?.distribution?.length) return;
    const chart = echarts.init(behaviorDistRef.current);
    chart.setOption({
      tooltip: { trigger: "item" },
      legend: { top: "bottom" },
      series: [{
        type: "pie",
        radius: ["35%", "65%"],
        avoidLabelOverlap: false,
        label: { show: true, formatter: "{b}\n{d}%" },
        data: behaviorStats.distribution.map((d) => ({
          name: d.behavior_type,
          value: d.count,
        })),
      }],
    });
    return () => chart.dispose();
  }, [behaviorStats]);

  useEffect(() => {
    if (!dailyTrendRef.current || !behaviorStats?.daily_trend?.length) return;
    const chart = echarts.init(dailyTrendRef.current);
    chart.setOption({
      tooltip: { trigger: "axis" },
      xAxis: { type: "category", data: behaviorStats.daily_trend.map((d) => d.date.slice(5)) },
      yAxis: { type: "value" },
      series: [{
        data: behaviorStats.daily_trend.map((d) => d.count),
        type: "line",
        smooth: true,
        areaStyle: { opacity: 0.15 },
        lineStyle: { color: "#6366f1", width: 2 },
        itemStyle: { color: "#6366f1" },
      }],
      grid: { top: 20, left: 50, right: 20, bottom: 20, containLabel: true },
    });
    return () => chart.dispose();
  }, [behaviorStats]);

  useEffect(() => {
    if (!qualityDistRef.current || !contentPerf?.quality_distribution?.length) return;
    const chart = echarts.init(qualityDistRef.current);
    chart.setOption({
      tooltip: { trigger: "axis" },
      xAxis: { type: "category", data: contentPerf.quality_distribution.map((d) => d.range) },
      yAxis: { type: "value" },
      series: [{
        data: contentPerf.quality_distribution.map((d) => d.count),
        type: "bar",
        barWidth: "50%",
        itemStyle: { color: "#06b6d4", borderRadius: [4, 4, 0, 0] },
      }],
      grid: { top: 20, left: 50, right: 20, bottom: 20, containLabel: true },
    });
    return () => chart.dispose();
  }, [contentPerf]);

  useEffect(() => {
    if (!boardDistRef.current || !contentPerf?.content_count_by_board?.length) return;
    const chart = echarts.init(boardDistRef.current);
    chart.setOption({
      tooltip: { trigger: "axis" },
      xAxis: {
        type: "category",
        data: contentPerf.content_count_by_board.map((d) => d.board_name),
        axisLabel: { rotate: 30 },
      },
      yAxis: { type: "value" },
      series: [{
        data: contentPerf.content_count_by_board.map((d) => d.count),
        type: "bar",
        barWidth: "40%",
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "#a78bfa" },
            { offset: 1, color: "#6366f1" },
          ]),
          borderRadius: [4, 4, 0, 0],
        },
      }],
      grid: { top: 20, left: 50, right: 20, bottom: 50, containLabel: true },
    });
    return () => chart.dispose();
  }, [contentPerf]);

  useEffect(() => {
    if (!tagCloudRef.current || !userAnalysis?.top_interest_tags?.length) return;
    const chart = echarts.init(tagCloudRef.current);
    const tags = [...userAnalysis.top_interest_tags].slice(0, 10);
    chart.setOption({
      tooltip: { trigger: "axis", axisPointer: { type: "shadow" } },
      xAxis: {
        type: "value",
        name: "用户数",
      },
      yAxis: {
        type: "category",
        data: tags.map((t) => t.tag_name || `标签 #${t.tag_id}`).reverse(),
        axisLabel: { fontSize: 12 },
      },
      series: [{
        data: tags.map((t) => t.user_count).reverse(),
        type: "bar",
        barWidth: "60%",
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: "#a78bfa" },
            { offset: 1, color: "#f472b6" },
          ]),
          borderRadius: [0, 4, 4, 0],
        },
        label: { show: true, position: "right", fontSize: 11 },
      }],
      grid: { top: 10, left: 100, right: 40, bottom: 10, containLabel: true },
    });
    return () => chart.dispose();
  }, [userAnalysis]);

  const loading = isLoading || bhLoading || uaLoading || cpLoading || riskLoading;
  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* 概览指标卡片 */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-primary">
            <Activity className="w-6 h-6" />
          </div>
          <div className="stat-title">行为总数</div>
          <div className="stat-value text-2xl">
            {(overview?.total_behaviors || 0).toLocaleString()}
          </div>
          <div className="stat-desc">今日 +{overview?.today_behaviors?.toLocaleString() || 0}</div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-info">
            <Eye className="w-6 h-6" />
          </div>
          <div className="stat-title">今日曝光</div>
          <div className="stat-value text-2xl">
            {(overview?.today_impressions || 0).toLocaleString()}
          </div>
          <div className="stat-desc">
            点击率 {overview?.click_rate || 0}%
          </div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-success">
            <Users className="w-6 h-6" />
          </div>
          <div className="stat-title">追踪用户</div>
          <div className="stat-value text-2xl">
            {(overview?.user_count || 0).toLocaleString()}
          </div>
          <div className="stat-desc">
            {(overview?.total_feedbacks || 0).toLocaleString()} 条反馈
          </div>
        </div>

        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-warning">
            <Target className="w-6 h-6" />
          </div>
          <div className="stat-title">已分析内容</div>
          <div className="stat-value text-2xl">
            {(overview?.content_count || 0).toLocaleString()}
          </div>
          <div className="stat-desc">
            均分 {overview?.avg_quality_score?.toFixed(1) || "0.0"}
          </div>
        </div>
      </div>

      {/* 图表区域 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {/* 行为分布饼图 */}
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
              <MousePointerClick className="w-4 h-4 text-primary" />
              行为分布
            </h4>
            <div ref={behaviorDistRef} className="h-72 w-full" />
          </div>
        </div>

        {/* 每日行为趋势 */}
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
              <TrendingUp className="w-4 h-4 text-info" />
              每日行为趋势 (7天)
            </h4>
            <div ref={dailyTrendRef} className="h-72 w-full" />
          </div>
        </div>

        {/* 质量分分布 */}
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
              <Target className="w-4 h-4 text-success" />
              内容质量分分布
            </h4>
            <div ref={qualityDistRef} className="h-72 w-full" />
          </div>
        </div>

        {/* 板块内容分布 */}
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
              <Activity className="w-4 h-4 text-secondary" />
              板块内容分布
            </h4>
            <div ref={boardDistRef} className="h-72 w-full" />
          </div>
        </div>
      </div>

      {/* 活跃用户排行 */}
      {userAnalysis?.top_active_users && userAnalysis.top_active_users.length > 0 && (
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-3">
              <Users className="w-4 h-4 text-primary" />
              最活跃用户 (近7天)
              <span className="text-base-content/40 text-xs ml-2">
                人均 {(userAnalysis.avg_behaviors_per_user || 0)} 次行为 | 今日活跃 {userAnalysis.active_users_today || 0} 人
              </span>
            </h4>
            <div className="overflow-x-auto">
              <table className="table table-sm">
                <thead>
                  <tr>
                    <th>#</th>
                    <th>用户</th>
                    <th>行为次数</th>
                    <th>主要行为</th>
                  </tr>
                </thead>
                <tbody>
                  {userAnalysis.top_active_users.map((u, i) => (
                    <tr key={u.user_id}>
                      <td className="font-bold">{i + 1}</td>
                      <td>{u.nickname || u.username}</td>
                      <td>{u.behavior_count.toLocaleString()}</td>
                      <td>
                        <span className="badge badge-sm badge-ghost">{u.top_behavior}</span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* 热门兴趣标签 */}
      {userAnalysis?.top_interest_tags && userAnalysis.top_interest_tags.length > 0 && (
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-3">
              <TrendingUp className="w-4 h-4 text-accent" />
              热门兴趣标签
            </h4>
            <div className="flex flex-wrap gap-3">
              {userAnalysis.top_interest_tags.map((tag) => (
                <div
                  key={tag.tag_id}
                  className="flex items-center gap-2 bg-base-200 rounded-lg px-3 py-2"
                >
                  <span className="text-sm font-medium">
                    {tag.tag_name || `标签 #${tag.tag_id}`}
                  </span>
                  <span className="text-xs text-base-content/50">
                    {tag.user_count} 人
                  </span>
                  <div
                    className="h-1.5 rounded-full bg-primary/30"
                    style={{
                      width: `${Math.min((tag.weight / (userAnalysis.top_interest_tags[0]?.weight || 1)) * 60, 60)}px`,
                    }}
                  />
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* 内容表现排行 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {contentPerf?.top_hot_content && contentPerf.top_hot_content.length > 0 && (
          <div className="card bg-base-100 border border-base-300 shadow-sm">
            <div className="card-body p-4">
              <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
                热度最高内容
              </h4>
              <div className="space-y-2">
                {contentPerf.top_hot_content.map((item, i) => (
                  <div key={item.creation_id} className="flex items-center gap-2 text-sm">
                    <span className="font-bold w-6">{i + 1}</span>
                    <span className="flex-1 truncate">{item.title || `内容 #${item.creation_id}`}</span>
                    <span className="text-base-content/50">{item.view_count} 浏览</span>
                    <span className="badge badge-sm">{item.hot_score?.toFixed(1)}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {contentPerf?.top_quality_content && contentPerf.top_quality_content.length > 0 && (
          <div className="card bg-base-100 border border-base-300 shadow-sm">
            <div className="card-body p-4">
              <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
                质量最高内容
              </h4>
              <div className="space-y-2">
                {contentPerf.top_quality_content.map((item, i) => (
                  <div key={item.creation_id} className="flex items-center gap-2 text-sm">
                    <span className="font-bold w-6">{i + 1}</span>
                    <span className="flex-1 truncate">{item.title || `内容 #${item.creation_id}`}</span>
                    <span className="text-base-content/50">{item.like_count} 赞</span>
                    <span className="badge badge-sm badge-success">{item.quality_score?.toFixed(1)}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>

      {/* 最热目标内容 */}
      {behaviorStats?.top_targets && behaviorStats.top_targets.length > 0 && (
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
              <Target className="w-4 h-4 text-error" />
              互动最多的内容
            </h4>
            <div className="space-y-1.5">
              {behaviorStats.top_targets.map((item, i) => (
                <div key={item.target_id} className="flex items-center gap-2 text-sm py-1">
                  <span className="font-bold w-6">{i + 1}</span>
                  <span className="flex-1 truncate">
                    {item.title || `内容 #${item.target_id}`}
                  </span>
                  <span className="badge badge-sm">{item.behavior_count} 次互动</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* 风控分析 */}
      {riskAnalysis && (
        <div className="space-y-4">
          <h3 className="text-lg font-bold flex items-center gap-2 mt-4">
            <Shield className="w-5 h-5 text-warning" />
            风控关联分析
          </h3>

          {/* 风控指标卡片 */}
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">
            <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm p-3">
              <div className="stat-figure text-warning"><AlertTriangle className="w-5 h-5" /></div>
              <div className="stat-title text-xs">风险用户</div>
              <div className="stat-value text-lg">{riskAnalysis.total_risk_users}</div>
              <div className="stat-desc text-xs text-error">高危 {riskAnalysis.danger_level_users}</div>
            </div>
            <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm p-3">
              <div className="stat-figure text-error"><Flag className="w-5 h-5" /></div>
              <div className="stat-title text-xs">违规总数</div>
              <div className="stat-value text-lg">{riskAnalysis.total_violations}</div>
              <div className="stat-desc text-xs">待处理 {riskAnalysis.pending_violations}</div>
            </div>
            <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm p-3">
              <div className="stat-figure text-error"><Ban className="w-5 h-5" /></div>
              <div className="stat-title text-xs">已封禁</div>
              <div className="stat-value text-lg">{riskAnalysis.total_bans}</div>
            </div>
            <div className="stat col-span-2 bg-base-100 rounded-lg border border-base-300 shadow-sm p-3">
              <div className="stat-title text-xs mb-2">违规类型分布</div>
              <div className="flex flex-wrap gap-1.5">
                {riskAnalysis.violation_distribution?.map((v) => (
                  <div key={v.violation_type} className="flex items-center gap-1 text-xs">
                    <span className="badge badge-xs badge-outline">{v.violation_type}</span>
                    <span className="text-base-content/50">{v.count}({v.ratio}%)</span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* 风控用户行为关联表 */}
          {riskAnalysis.risk_user_behaviors?.length > 0 && (
            <div className="card bg-base-100 border border-base-300 shadow-sm">
              <div className="card-body p-4">
                <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
                  <Shield className="w-4 h-4 text-warning" />
                  风控用户行为关联 TOP10
                </h4>
                <div className="overflow-x-auto">
                  <table className="table table-sm">
                    <thead>
                      <tr>
                        <th>用户</th>
                        <th>风险等级</th>
                        <th>违规次数</th>
                        <th>行为次数</th>
                        <th>主要行为</th>
                      </tr>
                    </thead>
                    <tbody>
                      {riskAnalysis.risk_user_behaviors.map((u) => (
                        <tr key={u.user_id}>
                          <td>{u.username || `用户 #${u.user_id}`}</td>
                          <td>
                            <span className={`badge badge-xs ${u.risk_level === "danger" ? "badge-error" : "badge-warning"}`}>
                              {u.risk_level}
                            </span>
                          </td>
                          <td>{u.violation_count}</td>
                          <td>{u.behavior_count || 0}</td>
                          <td><span className="badge badge-xs badge-ghost">{u.top_behavior || "-"}</span></td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            </div>
          )}

          {/* 被举报内容排行 */}
          {riskAnalysis.top_reported_content?.length > 0 && (
            <div className="card bg-base-100 border border-base-300 shadow-sm">
              <div className="card-body p-4">
                <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
                  <Flag className="w-4 h-4 text-error" />
                  被举报最多的内容
                </h4>
                <div className="space-y-1.5">
                  {riskAnalysis.top_reported_content.map((item, i) => (
                    <div key={item.creation_id} className="flex items-center gap-2 text-sm py-1">
                      <span className="font-bold w-6">{i + 1}</span>
                      <span className="flex-1 truncate">{item.title || `内容 #${item.creation_id}`}</span>
                      <span className={`badge badge-xs ${item.risk_level === "danger" ? "badge-error" : "badge-warning"}`}>
                        {item.report_count} 次举报
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
