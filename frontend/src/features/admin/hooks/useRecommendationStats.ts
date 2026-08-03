import { useQuery } from "@tanstack/react-query";
import { adminKeys } from "./useAdminKeys";
import {
  adminRecApi,
  RecOverviewStats,
  BehaviorStats,
  UserAnalysis,
  ContentPerformance,
  RiskAnalysis,
  ComprehensiveUserAnalysis,
} from "@/shared/api/modules/admin/recommendation";

export const recAdminKeys = {
  all: [...adminKeys.all, "recommendation"] as const,
  overview: () => [...recAdminKeys.all, "overview"] as const,
  behaviors: (days?: number) => [...recAdminKeys.all, "behaviors", days] as const,
  users: (days?: number) => [...recAdminKeys.all, "users", days] as const,
  content: () => [...recAdminKeys.all, "content"] as const,
  risk: () => [...recAdminKeys.all, "risk"] as const,
  userAnalysis: () => [...recAdminKeys.all, "user-analysis"] as const,
};

export function useRecOverview(enabled = true) {
  return useQuery<RecOverviewStats>({
    queryKey: recAdminKeys.overview(),
    queryFn: async () => {
      const res = await adminRecApi.getOverview();
      if (res.data.code !== 0) throw new Error(res.data.message || "获取概览失败");
      return res.data.data as RecOverviewStats;
    },
    enabled,
    staleTime: 60_000,
  });
}

export function useRecBehaviorStats(days = 7, enabled = true) {
  return useQuery<BehaviorStats>({
    queryKey: recAdminKeys.behaviors(days),
    queryFn: async () => {
      const res = await adminRecApi.getBehaviorStats(days);
      if (res.data.code !== 0) throw new Error(res.data.message || "获取行为统计失败");
      return res.data.data as BehaviorStats;
    },
    enabled,
    staleTime: 60_000,
  });
}

export function useRecUserAnalysis(days = 7, enabled = true) {
  return useQuery<UserAnalysis>({
    queryKey: recAdminKeys.users(days),
    queryFn: async () => {
      const res = await adminRecApi.getUserAnalysis(days);
      if (res.data.code !== 0) throw new Error(res.data.message || "获取用户分析失败");
      return res.data.data as UserAnalysis;
    },
    enabled,
    staleTime: 60_000,
  });
}

export function useRecContentPerformance(enabled = true) {
  return useQuery<ContentPerformance>({
    queryKey: recAdminKeys.content(),
    queryFn: async () => {
      const res = await adminRecApi.getContentPerformance();
      if (res.data.code !== 0) throw new Error(res.data.message || "获取内容表现失败");
      return res.data.data as ContentPerformance;
    },
    enabled,
    staleTime: 120_000,
  });
}

export function useRecRiskAnalysis(enabled = true) {
  return useQuery<RiskAnalysis>({
    queryKey: recAdminKeys.risk(),
    queryFn: async () => {
      const res = await adminRecApi.getRiskAnalysis();
      if (res.data.code !== 0) throw new Error(res.data.message || "获取风险分析失败");
      return res.data.data as RiskAnalysis;
    },
    enabled,
    staleTime: 60_000,
  });
}

export function useComprehensiveUserAnalysis(enabled = true) {
  return useQuery<ComprehensiveUserAnalysis>({
    queryKey: recAdminKeys.userAnalysis(),
    queryFn: async () => {
      const res = await adminRecApi.getComprehensiveUserAnalysis();
      if (res.data.code !== 0) throw new Error(res.data.message || "获取用户分析失败");
      return res.data.data as ComprehensiveUserAnalysis;
    },
    enabled,
    staleTime: 60_000,
  });
}
