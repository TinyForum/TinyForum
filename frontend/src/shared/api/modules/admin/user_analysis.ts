import apiClient from "../../client";
import { ApiResponse } from "../../types/basic.model";

// ── Types ──

export interface UATrendPoint {
  date: string;
  dau: number;
  new_user: number;
  actions: number;
}
export interface UAOverview {
  total_users: number;
  dau: number;
  wau: number;
  mau: number;
  stickiness: number;
  new_users_today: number;
  new_users_week: number;
  retention_day7: number;
  retention_day30: number;
  avg_daily_actions: number;
  avg_session_duration: number;
  trend_points: UATrendPoint[];
}

export interface UAProfile {
  user_id: number; username: string; nickname: string; avatar: string; email: string;
  joined_at: string; last_active_at: string; engagement_tier: string;
  total_actions: number; active_days: number; active_tags: string[];
  top_behaviors: Record<string, number>; content_count: number;
  comment_count: number; like_received: number;
  risk_level: string; violation_count: number; is_banned: boolean;
}
export interface UAProfileList {
  profiles: UAProfile[]; total: number; page: number; page_size: number;
}

export interface UASegment {
  segment_name: string; user_count: number; percentage: number;
  avg_actions: number; avg_active_days: number; retention_rate: number;
  top_tags: string[]; risk_ratio: number;
}
export interface UASegments { segments: UASegment[]; total_users: number; }

export interface UAFunnelStep { step_name: string; user_count: number; conversion: number; drop_off: number; }
export interface UABehaviorFunnel { steps: UAFunnelStep[]; }
export interface UAEventItem { event_type: string; count: number; unique_users: number; ratio: number; }
export interface UAHourlyItem { hour: number; count: number; }
export interface UAEventUserItem { user_id: number; username: string; count: number; event_type: string; }
export interface UABehavior {
  funnel: UABehaviorFunnel; event_distribution: UAEventItem[];
  hourly_heatmap: UAHourlyItem[]; top_event_users: UAEventUserItem[];
}

export interface UACohort {
  cohort_label: string; initial_users: number;
  retention_w1: number; retention_w2: number; retention_w3: number; retention_w4: number; retention_w8: number;
}
export interface UACohorts { cohorts: UACohort[]; week_labels: string[]; max_weeks: number; }

export interface UARiskScoreDist { range: string; count: number; }
export interface UARiskyUserItem {
  user_id: number; username: string; risk_level: string; violation_count: number;
  behavior_count: number; last_flag_reason: string; is_banned: boolean;
}
export interface UAFlaggedItem {
  creation_id: number; title: string; author_name: string;
  report_count: number; flag_reason: string; created_at: string;
}
export interface UARisk {
  total_risk_users: number; high_risk_count: number; new_risks_today: number;
  pending_reviews: number; score_distribution: UARiskScoreDist[];
  violation_trend: UATrendPoint[]; top_risky_users: UARiskyUserItem[];
  flagged_content_queue: UAFlaggedItem[];
}

// ── API ──

export const uaApi = {
  overview: () => apiClient.get<ApiResponse<UAOverview>>("/admin/recommendations/ua/overview"),
  segments: () => apiClient.get<ApiResponse<UASegments>>("/admin/recommendations/ua/segments"),
  profiles: (params?: { page?: number; page_size?: number; keyword?: string; tier?: string; sort_by?: string }) =>
    apiClient.get<ApiResponse<UAProfileList>>("/admin/recommendations/ua/profiles", { params }),
  behavior: () => apiClient.get<ApiResponse<UABehavior>>("/admin/recommendations/ua/behavior"),
  cohorts: () => apiClient.get<ApiResponse<UACohorts>>("/admin/recommendations/ua/cohorts"),
  risk: () => apiClient.get<ApiResponse<UARisk>>("/admin/recommendations/ua/risk"),
};
