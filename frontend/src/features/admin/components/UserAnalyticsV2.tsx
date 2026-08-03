"use client";

import { useState, useEffect, useRef } from "react";
import * as echarts from "echarts";
import {
  Users, Activity, PieChart, TrendingUp, Shield, Layers, BarChart3,
  AlertTriangle, Ban, Clock, Target, UserX,
} from "lucide-react";
import {
  useUAOverview, useUASegments, useUAProfiles, useUABehavior, useUACohorts, useUARisk,
} from "../hooks/useUAAdmin";
import { UATrendPoint, UAHourlyItem, UACohort } from "@/shared/api/modules/admin/user_analysis";

const tabs = [
  { id: "overview", label: "概览", icon: BarChart3 },
  { id: "segments", label: "用户分群", icon: Layers },
  { id: "profiles", label: "用户画像", icon: Users },
  { id: "behavior", label: "行为分析", icon: Activity },
  { id: "cohorts", label: "同期群分析", icon: Clock },
  { id: "risk", label: "风险评估", icon: Shield },
];

function TrendChart({ data }: { data: UATrendPoint[] }) {
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!ref.current || !data?.length) return;
    const c = echarts.init(ref.current);
    c.setOption({
      tooltip: { trigger: "axis" },
      legend: { data: ["DAU", "New Users", "Actions"] },
      xAxis: { type: "category", data: data.map((d) => d.date.slice(5)) },
      yAxis: { type: "value" },
      series: [
        { name: "DAU", data: data.map((d) => d.dau), type: "line", smooth: true, lineStyle: { width: 2, color: "#6366f1" } },
        { name: "New Users", data: data.map((d) => d.new_user), type: "line", smooth: true, lineStyle: { width: 2, color: "#22c55e" } },
        { name: "Actions", data: data.map((d) => d.actions), type: "bar", barWidth: "30%", itemStyle: { color: "#f59e0b", opacity: 0.6 } },
      ],
      grid: { top: 40, left: 60, right: 20, bottom: 30 },
    });
    return () => c.dispose();
  }, [data]);
  return <div ref={ref} className="h-72 w-full" />;
}

function HourlyHeatmap({ data }: { data: UAHourlyItem[] }) {
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!ref.current || !data?.length) return;
    const c = echarts.init(ref.current);
    c.setOption({
      tooltip: { trigger: "axis" },
      xAxis: { type: "category", data: data.map((d) => d.hour + ":00"), axisLabel: { fontSize: 10 } },
      yAxis: { type: "value" },
      series: [{
        data: data.map((d) => d.count), type: "bar", barWidth: "70%",
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: "#a78bfa" }, { offset: 1, color: "#6366f1" },
          ]),
          borderRadius: [2, 2, 0, 0],
        },
      }],
      grid: { top: 10, left: 50, right: 20, bottom: 30 },
    });
    return () => c.dispose();
  }, [data]);
  return <div ref={ref} className="h-48 w-full" />;
}

function RetentionHeatmap({ cohorts }: { cohorts: UACohort[] }) {
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!ref.current || !cohorts?.length) return;
    const c = echarts.init(ref.current);
    const weeks = ["W1", "W2", "W3", "W4", "W8"];
    const data = cohorts.slice(0, 7).map((co) => [
      co.retention_w1 || 0, co.retention_w2 || 0, co.retention_w3 || 0, co.retention_w4 || 0, co.retention_w8 || 0,
    ]);
    const labels = cohorts.slice(0, 7).map((co) => co.cohort_label);
    const chartData: { value: [number, number, number]; label?: { show: boolean; formatter: string } }[] = [];
    data.forEach((row, ri) => {
      row.forEach((val, ci) => {
        chartData.push({ value: [ci, ri, val], label: val > 0 ? { show: true, formatter: val.toFixed(1) + "%" } : { show: false, formatter: "" } });
      });
    });
    c.setOption({
      tooltip: {},
      xAxis: { type: "category", data: weeks, splitArea: { show: true } },
      yAxis: { type: "category", data: labels.reverse(), splitArea: { show: true } },
      visualMap: { min: 0, max: 100, calculable: true, orient: "horizontal", left: "center", bottom: 0, inRange: { color: ["#e0e7ff", "#6366f1", "#312e81"] } },
      series: [{
        type: "heatmap", data: chartData,
        label: { show: true, fontSize: 10 },
        emphasis: { itemStyle: { shadowBlur: 10, shadowColor: "rgba(0,0,0,0.5)" } },
      }],
      grid: { top: 10, left: 110, right: 20, bottom: 60 },
    });
    return () => c.dispose();
  }, [cohorts]);
  return <div ref={ref} className="h-96 w-full" />;
}

function KPI({ icon: Icon, label, value, sub, color }: { icon: typeof Users; label: string; value: string; sub?: string; color: string }) {
  return (
    <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm p-3">
      <div className={`stat-figure ${color}`}><Icon className="w-5 h-5" /></div>
      <div className="stat-title text-xs">{label}</div>
      <div className="stat-value text-lg">{value}</div>
      {sub && <div className="stat-desc text-xs">{sub}</div>}
    </div>
  );
}

export function UserAnalyticsV2() {
  const [tab, setTab] = useState("overview");
  const [profileParams, setProfileParams] = useState({ page: 1, page_size: 20, keyword: "", tier: "", sort_by: "" });
  const { data: overview } = useUAOverview(tab === "overview");
  const { data: segments } = useUASegments(tab === "segments");
  const { data: profiles } = useUAProfiles(profileParams as unknown as Record<string, unknown>, tab === "profiles");
  const { data: behavior } = useUABehavior(tab === "behavior");
  const { data: cohorts } = useUACohorts(tab === "cohorts");
  const { data: risk } = useUARisk(tab === "risk");

  return (
    <div className="space-y-4">
      {/* Tab Bar */}
      <div className="flex border-b border-base-300 bg-base-100 rounded-t-lg overflow-x-auto">
        {tabs.map((t) => {
          const Icon = t.icon;
          return (
            <button key={t.id} onClick={() => setTab(t.id)}
              className={`flex items-center gap-2 px-4 py-3 text-sm font-medium whitespace-nowrap transition-colors ${tab === t.id ? "text-primary border-b-2 border-primary bg-primary/5" : "text-base-content/60 hover:text-base-content"}`}>
              <Icon className="w-4 h-4" />{t.label}
            </button>
          );
        })}
      </div>

      {/* ── Tab: Overview ── */}
      {tab === "overview" && overview && (
        <div className="space-y-4">
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3">
            <KPI icon={Users} label="DAU" value={overview.dau?.toLocaleString()} sub={`Stickiness ${overview.stickiness}%`} color="text-primary" />
            <KPI icon={Users} label="WAU" value={overview.wau?.toLocaleString()} color="text-primary" />
            <KPI icon={Users} label="MAU" value={overview.mau?.toLocaleString()} color="text-primary" />
            <KPI icon={TrendingUp} label="今日新增" value={overview.new_users_today?.toLocaleString()} sub={`7日: ${overview.new_users_week}`} color="text-success" />
            <KPI icon={Clock} label="7日留存" value={`${overview.retention_day7}%`} sub={`30日: ${overview.retention_day30}%`} color="text-info" />
            <KPI icon={Activity} label="人均日行为" value={overview.avg_daily_actions?.toFixed(1)} color="text-accent" />
          </div>
          <div className="card bg-base-100 border border-base-300 shadow-sm p-4">
            <h4 className="text-sm font-medium mb-2 flex items-center gap-1"><TrendingUp className="w-4 h-4 text-primary" />14日趋势</h4>
            <TrendChart data={overview.trend_points} />
          </div>
        </div>
      )}

      {/* ── Tab: Segments ── */}
      {tab === "segments" && segments && (
        <div className="space-y-4">
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">
            {segments.segments.map((seg) => (
              <div key={seg.segment_name} className="card bg-base-100 border border-base-300 shadow-sm p-3">
                <h4 className="font-semibold text-sm capitalize">{seg.segment_name}</h4>
                <div className="text-2xl font-bold">{seg.user_count.toLocaleString()}</div>
                <div className="text-xs text-base-content/50">{seg.percentage}% of total</div>
                <div className="mt-2 text-xs space-y-1">
                  <div className="flex justify-between"><span>Avg Actions</span><span className="font-mono">{seg.avg_actions?.toFixed(1)}</span></div>
                  <div className="flex justify-between"><span>Risk Ratio</span><span className={`font-mono ${(seg.risk_ratio || 0) > 0.1 ? "text-error" : ""}`}>{((seg.risk_ratio || 0) * 100).toFixed(1)}%</span></div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* ── Tab: Profiles ── */}
      {tab === "profiles" && (
        <div className="space-y-4">
          <div className="flex flex-wrap gap-2 items-center">
            <div className="join">
              <input placeholder="搜索用户..." value={profileParams.keyword} onChange={(e) => setProfileParams((p) => ({ ...p, keyword: e.target.value, page: 1 }))}
                className="input input-bordered input-sm join-item w-48" />
            </div>
            <select value={profileParams.tier} onChange={(e) => setProfileParams((p) => ({ ...p, tier: e.target.value, page: 1 }))}
              className="select select-bordered select-sm">
              <option value="">全部分群</option>
              <option value="power">Power</option>
              <option value="core">Core</option>
              <option value="casual">Casual</option>
              <option value="dormant">Dormant</option>
              <option value="churned">Churned</option>
              <option value="risky">风险用户</option>
            </select>
            <select value={profileParams.sort_by} onChange={(e) => setProfileParams((p) => ({ ...p, sort_by: e.target.value }))}
              className="select select-bordered select-sm">
              <option value="">默认排序</option>
              <option value="actions">行为次数</option>
              <option value="active">最近活跃</option>
              <option value="risk">风险排序</option>
            </select>
            <span className="text-xs text-base-content/50 ml-auto">共 {profiles?.total || 0} 人</span>
          </div>
          {profiles?.profiles && (
            <div className="overflow-x-auto">
              <table className="table table-sm table-zebra">
                <thead>
                  <tr>
                    <th>用户</th><th>分群</th><th>活跃天数</th><th>行为次数</th>
                    <th>最近活跃</th><th>风险</th><th>违规</th><th>状态</th>
                  </tr>
                </thead>
                <tbody>
                  {profiles.profiles.map((u) => (
                    <tr key={u.user_id}>
                      <td>
                        <div className="flex items-center gap-2">
                          <span className="font-medium text-sm">{u.nickname || u.username}</span>
                          <span className="text-xs text-base-content/40">{u.email}</span>
                        </div>
                      </td>
                      <td><span className="badge badge-xs badge-outline capitalize">{u.engagement_tier}</span></td>
                      <td>{u.active_days}d</td>
                      <td>{u.total_actions?.toLocaleString()}</td>
                      <td className="text-xs">{u.last_active_at?.slice(0, 10)}</td>
                      <td><span className={`badge badge-xs ${u.risk_level !== "normal" ? "badge-warning" : "badge-ghost"}`}>{u.risk_level}</span></td>
                      <td className={u.violation_count > 0 ? "text-error font-bold" : ""}>{u.violation_count}</td>
                      <td>{u.is_banned ? <Ban className="w-4 h-4 text-error" /> : "-"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
          <div className="flex justify-center gap-2">
            <button onClick={() => setProfileParams((p) => ({ ...p, page: Math.max(1, p.page - 1) }))} className="btn btn-sm" disabled={profileParams.page <= 1}>上一页</button>
            <span className="text-sm self-center">第 {profileParams.page} 页</span>
            <button onClick={() => setProfileParams((p) => ({ ...p, page: p.page + 1 }))} className="btn btn-sm"
              disabled={!profiles || !profiles.profiles || profiles.profiles.length < profileParams.page_size}>下一页</button>
          </div>
        </div>
      )}

      {/* ── Tab: Behavior ── */}
      {tab === "behavior" && behavior && (
        <div className="space-y-4">
          {/* Funnel */}
          <div className="card bg-base-100 border border-base-300 shadow-sm p-4">
            <h4 className="text-sm font-medium mb-2 flex items-center gap-1"><Target className="w-4 h-4 text-primary" />行为漏斗（7天）</h4>
            <div className="flex items-end gap-3">
              {behavior.funnel.steps.map((s, i) => (
                <div key={s.step_name} className="flex-1 text-center">
                  <div className={`rounded-t-lg pt-3 pb-1 ${i === 0 ? "bg-primary/20" : i === 1 ? "bg-info/20" : i === 2 ? "bg-success/20" : "bg-warning/20"}`}
                    style={{ height: `${Math.max(8, s.user_count / Math.max(1, behavior.funnel.steps[0]?.user_count || 1) * 80)}px` }}>
                    <div className="font-bold text-sm">{s.user_count.toLocaleString()}</div>
                  </div>
                  <div className="bg-base-200 rounded-b-lg py-1 text-xs">{s.step_name}</div>
                  <div className="text-xs text-base-content/50 mt-1">{s.conversion}%</div>
                </div>
              ))}
            </div>
          </div>

          {/* Event Distribution */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div className="card bg-base-100 border border-base-300 shadow-sm p-4">
              <h4 className="text-sm font-medium mb-2 flex items-center gap-1"><PieChart className="w-4 h-4 text-info" />事件类型分布</h4>
              <table className="table table-xs">
                <thead><tr><th>类型</th><th>次数</th><th>独立用户</th><th>占比</th></tr></thead>
                <tbody>
                  {behavior.event_distribution?.map((e) => (
                    <tr key={e.event_type}><td className="font-mono text-xs">{e.event_type}</td><td>{e.count?.toLocaleString()}</td><td>{e.unique_users?.toLocaleString()}</td><td>{e.ratio}%</td></tr>
                  ))}
                </tbody>
              </table>
            </div>
            <div className="card bg-base-100 border border-base-300 shadow-sm p-4">
              <h4 className="text-sm font-medium mb-2 flex items-center gap-1"><Clock className="w-4 h-4 text-accent" />24小时活跃热度</h4>
              <HourlyHeatmap data={behavior.hourly_heatmap} />
            </div>
          </div>

          {/* Top Event Users */}
          {behavior.top_event_users?.length > 0 && (
            <div className="card bg-base-100 border border-base-300 shadow-sm p-4">
              <h4 className="text-sm font-medium mb-2">Top 活跃用户</h4>
              <div className="flex flex-wrap gap-2">
                {behavior.top_event_users.map((u, i) => (
                  <div key={u.user_id} className="bg-base-200 rounded-lg px-3 py-2 flex items-center gap-2">
                    <span className="font-bold text-sm">{i + 1}</span>
                    <span className="text-sm">{u.username}</span>
                    <span className="badge badge-xs">{u.event_type}</span>
                    <span className="text-xs text-base-content/50">{u.count}</span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* ── Tab: Cohorts ── */}
      {tab === "cohorts" && cohorts && (
        <div className="space-y-4">
          <div className="card bg-base-100 border border-base-300 shadow-sm p-4">
            <h4 className="text-sm font-medium mb-2 flex items-center gap-1"><Clock className="w-4 h-4 text-primary" />同期群留存热力图</h4>
            <RetentionHeatmap cohorts={cohorts.cohorts} />
          </div>
          <div className="overflow-x-auto">
            <table className="table table-sm table-zebra">
              <thead>
                <tr>
                  <th>Cohort</th><th>Users</th><th>W1</th><th>W2</th><th>W3</th><th>W4</th><th>W8</th>
                </tr>
              </thead>
              <tbody>
                {cohorts.cohorts?.filter((c) => c.initial_users > 0).map((c) => (
                  <tr key={c.cohort_label}>
                    <td className="font-mono text-xs">{c.cohort_label}</td>
                    <td>{c.initial_users}</td>
                    <td className={c.retention_w1 > 30 ? "text-success" : ""}>{c.retention_w1}%</td>
                    <td>{c.retention_w2}%</td>
                    <td>{c.retention_w3}%</td>
                    <td>{c.retention_w4}%</td>
                    <td>{c.retention_w8}%</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* ── Tab: Risk ── */}
      {tab === "risk" && risk && (
        <div className="space-y-4">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <KPI icon={AlertTriangle} label="风险用户" value={risk.total_risk_users?.toLocaleString()} sub={`高危: ${risk.high_risk_count}`} color="text-error" />
            <KPI icon={Clock} label="今日新风险" value={risk.new_risks_today?.toLocaleString()} color="text-warning" />
            <KPI icon={UserX} label="待审核" value={risk.pending_reviews?.toLocaleString()} color="text-error" />
            <KPI icon={Shield} label="已标记" value={risk.flagged_content_queue?.length?.toString() || "0"} sub="条内容待审" color="text-warning" />
          </div>

          {/* Violation Trend + Flagged Queue */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            <div className="card bg-base-100 border border-base-300 shadow-sm p-4">
              <h4 className="text-sm font-medium mb-2">违规趋势（14天）</h4>
              <TrendChart data={risk.violation_trend?.map((d) => ({ ...d, dau: 0, new_user: 0, actions: d.actions || 0 })) || []} />
            </div>
            <div className="card bg-base-100 border border-base-300 shadow-sm p-4">
              <h4 className="text-sm font-medium mb-2 flex items-center gap-1"><Target className="w-4 h-4 text-error" />被举报内容队列</h4>
              <div className="space-y-2 max-h-72 overflow-y-auto">
                {risk.flagged_content_queue?.map((f) => (
                  <div key={f.creation_id} className="flex items-center justify-between text-xs p-2 bg-base-200 rounded">
                    <div className="flex-1 truncate mr-2">
                      <span className="font-medium">{f.title || `#${f.creation_id}`}</span>
                      <span className="text-base-content/40 ml-2">by {f.author_name}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <span className="badge badge-xs badge-error">×{f.report_count}</span>
                      <span className="text-base-content/40">{f.created_at?.slice(0, 10)}</span>
                    </div>
                  </div>
                ))}
                {!risk.flagged_content_queue?.length && <span className="text-base-content/30">暂无被举报内容</span>}
              </div>
            </div>
          </div>

          {/* Top Risky Users */}
          {risk.top_risky_users?.length > 0 && (
            <div className="card bg-base-100 border border-base-300 shadow-sm p-4">
              <h4 className="text-sm font-medium mb-2">高风险用户 TOP15</h4>
              <div className="overflow-x-auto">
                <table className="table table-sm table-zebra">
                  <thead>
                    <tr><th>#</th><th>用户</th><th>风险等级</th><th>违规</th><th>行为</th><th>原因</th><th>状态</th></tr>
                  </thead>
                  <tbody>
                    {risk.top_risky_users.map((u, i) => (
                      <tr key={u.user_id}>
                        <td className="font-bold">{i + 1}</td>
                        <td>{u.username}</td>
                        <td><span className={`badge badge-xs ${u.risk_level === "danger" ? "badge-error" : "badge-warning"}`}>{u.risk_level}</span></td>
                        <td className="text-error font-bold">{u.violation_count}</td>
                        <td>{u.behavior_count?.toLocaleString()}</td>
                        <td className="text-xs max-w-[120px] truncate">{u.last_flag_reason}</td>
                        <td>{u.is_banned ? <Ban className="w-4 h-4 text-error" /> : "-"}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
