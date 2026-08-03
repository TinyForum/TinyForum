import { useQuery } from "@tanstack/react-query";
import { adminKeys } from "./useAdminKeys";
import { uaApi, UAOverview, UASegments, UAProfileList, UABehavior, UACohorts, UARisk } from "@/shared/api/modules/admin/user_analysis";

export const uaKeys = {
  all: [...adminKeys.all, "ua"] as const,
  overview: () => [...uaKeys.all, "overview"] as const,
  segments: () => [...uaKeys.all, "segments"] as const,
  profiles: (params?: Record<string, unknown>) => [...uaKeys.all, "profiles", params] as const,
  behavior: () => [...uaKeys.all, "behavior"] as const,
  cohorts: () => [...uaKeys.all, "cohorts"] as const,
  risk: () => [...uaKeys.all, "risk"] as const,
};

export function useUAOverview(enabled = true) {
  return useQuery({ queryKey: uaKeys.overview(), queryFn: async () => { const r = await uaApi.overview(); return r.data.data as UAOverview; }, enabled, staleTime: 30_000 });
}
export function useUASegments(enabled = true) {
  return useQuery({ queryKey: uaKeys.segments(), queryFn: async () => { const r = await uaApi.segments(); return r.data.data as UASegments; }, enabled, staleTime: 60_000 });
}
export function useUAProfiles(params?: Record<string, unknown>, enabled = true) {
  return useQuery({ queryKey: uaKeys.profiles(params), queryFn: async () => { const r = await uaApi.profiles(params as { page?: number; page_size?: number; keyword?: string; tier?: string; sort_by?: string }); return r.data.data as UAProfileList; }, enabled, staleTime: 15_000 });
}
export function useUABehavior(enabled = true) {
  return useQuery({ queryKey: uaKeys.behavior(), queryFn: async () => { const r = await uaApi.behavior(); return r.data.data as UABehavior; }, enabled, staleTime: 30_000 });
}
export function useUACohorts(enabled = true) {
  return useQuery({ queryKey: uaKeys.cohorts(), queryFn: async () => { const r = await uaApi.cohorts(); return r.data.data as UACohorts; }, enabled, staleTime: 120_000 });
}
export function useUARisk(enabled = true) {
  return useQuery({ queryKey: uaKeys.risk(), queryFn: async () => { const r = await uaApi.risk(); return r.data.data as UARisk; }, enabled, staleTime: 30_000 });
}
