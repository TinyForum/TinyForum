"use client";

import { useEffect, useRef } from "react";
import * as echarts from "echarts";
import { useComprehensiveUserAnalysis } from "../hooks/useRecommendationStats";
import {
  Users,
  Activity,
  Tag,
  Target,
  Shield,
  AlertTriangle,
  Ban,
  Hash,
  PieChart,
  Network,
} from "lucide-react";
import { ComprehensiveUserAnalysis } from "@/shared/api/modules/admin/recommendation";
import Image from "next/image";

function TagDistChart({ data }: { data: ComprehensiveUserAnalysis["tag_distribution"] }) {
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!ref.current || !data?.length) return;
    const chart = echarts.init(ref.current);
    const tags = data.slice(0, 10);
    chart.setOption({
      tooltip: { trigger: "axis", axisPointer: { type: "shadow" } },
      xAxis: { type: "value", name: "用户数" },
      yAxis: {
        type: "category",
        data: tags.map((t) => t.tag_name).reverse(),
        axisLabel: { fontSize: 11 },
      },
      series: [{
        type: "bar",
        data: tags.map((t) => t.user_count).reverse(),
        barWidth: "55%",
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: "#22c55e" },
            { offset: 1, color: "#06b6d4" },
          ]),
          borderRadius: [0, 4, 4, 0],
        },
        label: { show: true, position: "right", fontSize: 10, formatter: "{c}人" },
      }],
      grid: { top: 10, left: 90, right: 45, bottom: 10, containLabel: true },
    });
    return () => chart.dispose();
  }, [data]);
  return <div ref={ref} className="h-72 w-full" />;
}

function BehaviorTypeChart({ patterns }: { patterns: ComprehensiveUserAnalysis["user_behavior_patterns"] }) {
  const ref = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!ref.current || !patterns?.length) return;
    const chart = echarts.init(ref.current);
    const typeCount: Record<string, number> = {};
    patterns.forEach((p) => {
      Object.entries(p.behavior_breakdown || {}).forEach(([k, v]) => {
        typeCount[k] = (typeCount[k] || 0) + v;
      });
    });
    const data = Object.entries(typeCount)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 8);
    chart.setOption({
      tooltip: { trigger: "item", formatter: "{b}: {c} ({d}%)" },
      series: [{
        type: "pie",
        radius: ["40%", "70%"],
        label: { formatter: "{b}\n{d}%" },
        data: data.map(([name, value]) => ({ name, value })),
      }],
    });
    return () => chart.dispose();
  }, [patterns]);
  return <div ref={ref} className="h-64 w-full" />;
}

export function UserAnalytics() {
  const { data, isLoading } = useComprehensiveUserAnalysis(true);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg" />
      </div>
    );
  }

  if (!data) return <div className="text-center text-base-content/50 py-8">暂无数据</div>;

  const { overview } = data;

  return (
    <div className="space-y-6">
      {/* 概览指标 */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-primary"><Users className="w-5 h-5" /></div>
          <div className="stat-title text-xs">总用户</div>
          <div className="stat-value text-xl">{overview.total_users?.toLocaleString()}</div>
          <div className="stat-desc text-xs">{overview.shared_users_count?.toLocaleString()} 人有行为记录</div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-info"><Activity className="w-5 h-5" /></div>
          <div className="stat-title text-xs">行为记录</div>
          <div className="stat-value text-xl">{overview.total_behavior_records?.toLocaleString()}</div>
          <div className="stat-desc text-xs">24h 人均 {overview.avg_behaviors_per_user_24h}</div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-success"><Tag className="w-5 h-5" /></div>
          <div className="stat-title text-xs">有画像用户</div>
          <div className="stat-value text-xl">{overview.users_with_profile?.toLocaleString()}</div>
          <div className="stat-desc text-xs">人均 {overview.avg_tags_per_user?.toFixed(1)} 标签</div>
        </div>
        <div className="stat bg-base-100 rounded-lg border border-base-300 shadow-sm">
          <div className="stat-figure text-error"><AlertTriangle className="w-5 h-5" /></div>
          <div className="stat-title text-xs">风险用户</div>
          <div className="stat-value text-xl">{overview.risk_user_count}</div>
          <div className="stat-desc text-xs">违规 {overview.violation_user_count} 人</div>
        </div>
      </div>

      {/* 标签用户分布 + 行为类型分布 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
              <Hash className="w-4 h-4 text-success" /> 标签用户分布 TOP10
            </h4>
            <TagDistChart data={data.tag_distribution} />
          </div>
        </div>
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-2">
              <PieChart className="w-4 h-4 text-info" /> 全局行为类型分布
            </h4>
            <BehaviorTypeChart patterns={data.user_behavior_patterns} />
          </div>
        </div>
      </div>

      {/* 用户行为模式表 */}
      {data.user_behavior_patterns?.length > 0 && (
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-3">
              <Activity className="w-4 h-4 text-primary" />
              用户行为模式 TOP20
            </h4>
            <div className="overflow-x-auto">
              <table className="table table-sm table-zebra">
                <thead>
                  <tr>
                    <th>#</th><th>用户</th><th>行为次数</th>
                    <th>行为分布</th><th>兴趣标签</th>
                    <th>风险</th><th>违规</th>
                  </tr>
                </thead>
                <tbody>
                  {data.user_behavior_patterns.map((u, i) => (
                    <tr key={u.user_id}>
                      <td className="font-bold">{i + 1}</td>
                      <td>
                        <div className="flex items-center gap-2">
                          {u.avatar && <Image src={u.avatar} alt="" width={24} height={24} className="w-6 h-6 rounded-full" />}
                          <span>{u.nickname || u.username}</span>
                        </div>
                      </td>
                      <td>{u.total_behaviors?.toLocaleString()}</td>
                      <td>
                        <div className="flex flex-wrap gap-1">
                          {Object.entries(u.behavior_breakdown || {}).slice(0, 4).map(([k, v]) => (
                            <span key={k} className="badge badge-xs badge-ghost">{k}:{v}</span>
                          ))}
                        </div>
                      </td>
                      <td>
                        <div className="flex flex-wrap gap-1">
                          {(u.active_tags || []).slice(0, 3).map((t, j) => (
                            <span key={j} className="badge badge-xs badge-outline">{t}</span>
                          ))}
                        </div>
                      </td>
                      <td>
                        <span className={`badge badge-xs ${u.risk_level !== "normal" ? "badge-warning" : "badge-ghost"}`}>
                          {u.risk_level}
                        </span>
                      </td>
                      <td className={u.violation_count > 0 ? "text-error font-bold" : ""}>
                        {u.violation_count || 0}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* 相似用户群组 */}
      {data.similar_user_groups?.length > 0 && (
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-3">
              <Network className="w-4 h-4 text-secondary" />
              相似用户群组（基于行为协同过滤）
            </h4>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {data.similar_user_groups.map((group) => (
                <div key={group.seed_user_id} className="bg-base-200 rounded-lg p-3">
                  <div className="flex items-center gap-2 mb-2">
                    <Target className="w-4 h-4 text-primary" />
                    <span className="font-semibold text-sm">{group.seed_username}</span>
                    <span className="text-xs text-base-content/40">种子用户</span>
                  </div>
                  <div className="space-y-1">
                    {group.similar_users.map((su) => (
                      <div key={su.user_id} className="flex items-center gap-2 text-xs pl-4">
                        <span className="text-primary font-mono">{su.similarity_score}</span>
                        <span>{su.nickname || su.username}</span>
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* 用户风险画像 */}
      {data.risk_user_list?.length > 0 && (
        <div className="card bg-base-100 border border-base-300 shadow-sm">
          <div className="card-body p-4">
            <h4 className="text-sm font-medium flex items-center gap-1 mb-3">
              <Shield className="w-4 h-4 text-error" />
              用户风险画像 TOP15
            </h4>
            <div className="overflow-x-auto">
              <table className="table table-sm table-zebra">
                <thead>
                  <tr>
                    <th>#</th><th>用户</th><th>风险等级</th>
                    <th>违规次数</th><th>最近违规</th>
                    <th>行为次数</th><th>状态</th>
                  </tr>
                </thead>
                <tbody>
                  {data.risk_user_list.map((u, i) => (
                    <tr key={u.user_id}>
                      <td className="font-bold">{i + 1}</td>
                      <td>{u.username}</td>
                      <td>
                        <span className={`badge badge-xs ${u.risk_level === "danger" ? "badge-error" : "badge-warning"}`}>
                          {u.risk_level}
                        </span>
                      </td>
                      <td className="font-bold text-error">{u.violation_count}</td>
                      <td className="text-xs">
                        <span className="badge badge-xs badge-ghost">{u.last_violation_type}</span>
                        <span className="text-base-content/40 ml-1">{u.last_violation_at?.slice(0, 10)}</span>
                      </td>
                      <td>{u.behavior_count?.toLocaleString()}</td>
                      <td>
                        {u.is_banned ? (
                          <Ban className="w-4 h-4 text-error" />
                        ) : (
                          <span className="text-base-content/40">-</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
