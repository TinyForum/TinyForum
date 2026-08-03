import apiClient from "../../client";
import { ApiResponse } from "../../types/basic.model";

export interface RecOverviewStats {
  total_behaviors: number;
  total_feedbacks: number;
  today_behaviors: number;
  today_impressions: number;
  click_rate: number;
  dismiss_rate: number;
  user_count: number;
  content_count: number;
  avg_quality_score: number;
}

export interface BehaviorDistribution {
  behavior_type: string;
  count: number;
  ratio: number;
}

export interface TrendPoint {
  date: string;
  count: number;
}

export interface TopTargetItem {
  target_id: number;
  title: string;
  behavior_count: number;
}

export interface BehaviorStats {
  distribution: BehaviorDistribution[];
  daily_trend: TrendPoint[];
  top_targets: TopTargetItem[];
}

export interface ActiveUserItem {
  user_id: number;
  username: string;
  nickname: string;
  behavior_count: number;
  top_behavior: string;
}

export interface TagWeightItem {
  tag_id: number;
  tag_name: string;
  weight: number;
  user_count: number;
}

export interface UserAnalysis {
  total_tracked_users: number;
  active_users_today: number;
  avg_behaviors_per_user: number;
  top_active_users: ActiveUserItem[];
  top_interest_tags: TagWeightItem[];
}

export interface ContentPerfItem {
  creation_id: number;
  title: string;
  view_count: number;
  like_count: number;
  hot_score: number;
  quality_score: number;
}

export interface BoardContentCount {
  board_id: number;
  board_name: string;
  count: number;
}

export interface ScoreBucket {
  range: string;
  count: number;
}

export interface ContentPerformance {
  top_hot_content: ContentPerfItem[];
  top_quality_content: ContentPerfItem[];
  content_count_by_board: BoardContentCount[];
  quality_distribution: ScoreBucket[];
}

export const adminRecApi = {
  getOverview: () =>
    apiClient.get<ApiResponse<RecOverviewStats>>("/admin/recommendations/overview"),

  getBehaviorStats: (days?: number) =>
    apiClient.get<ApiResponse<BehaviorStats>>("/admin/recommendations/behaviors", {
      params: { days },
    }),

  getUserAnalysis: (days?: number) =>
    apiClient.get<ApiResponse<UserAnalysis>>("/admin/recommendations/users", {
      params: { days },
    }),

  getContentPerformance: () =>
    apiClient.get<ApiResponse<ContentPerformance>>("/admin/recommendations/content"),

  getRiskAnalysis: () =>
    apiClient.get<ApiResponse<RiskAnalysis>>("/admin/recommendations/risk-analysis"),

  getComprehensiveUserAnalysis: () =>
    apiClient.get<ApiResponse<ComprehensiveUserAnalysis>>("/admin/recommendations/user-analysis"),
};

export interface ViolationTypeItem {
  violation_type: string;
  count: number;
  ratio: number;
}

export interface RiskUserBehavior {
  user_id: number;
  username: string;
  risk_level: string;
  violation_count: number;
  behavior_count: number;
  top_behavior: string;
}

export interface ReportedContentItem {
  creation_id: number;
  title: string;
  report_count: number;
  risk_level: string;
}

export interface RiskAnalysis {
  total_risk_users: number;
  danger_level_users: number;
  total_violations: number;
  pending_violations: number;
  total_bans: number;
  violation_distribution: ViolationTypeItem[];
  risk_user_behaviors: RiskUserBehavior[];
  top_reported_content: ReportedContentItem[];
}

// --- 用户综合分析 ---
export interface UserAnalysisOverview {
  total_users: number;
  total_behavior_records: number;
  users_with_profile: number;
  avg_tags_per_user: number;
  avg_behaviors_per_user_24h: number;
  shared_users_count: number;
  risk_user_count: number;
  violation_user_count: number;
}

export interface TagUserDistribution {
  tag_id: number;
  tag_name: string;
  user_count: number;
  avg_weight: number;
  post_count: number;
}

export interface UserBehaviorPattern {
  user_id: number;
  username: string;
  nickname: string;
  avatar: string;
  total_behaviors: number;
  behavior_breakdown: Record<string, number>;
  active_tags: string[];
  last_active_at: string;
  risk_level: string;
  violation_count: number;
}

export interface SimilarUserItem {
  user_id: number;
  username: string;
  nickname: string;
  avatar: string;
  similarity_score: number;
  shared_tags: string[];
  common_behavior: string;
}

export interface SimilarUserGroup {
  seed_user_id: number;
  seed_username: string;
  similar_users: SimilarUserItem[];
}

export interface UserRiskProfile {
  user_id: number;
  username: string;
  risk_level: string;
  violation_count: number;
  behavior_count: number;
  last_violation_type: string;
  last_violation_at: string;
  is_banned: boolean;
}

export interface ComprehensiveUserAnalysis {
  overview: UserAnalysisOverview;
  tag_distribution: TagUserDistribution[];
  user_behavior_patterns: UserBehaviorPattern[];
  similar_user_groups: SimilarUserGroup[];
  risk_user_list: UserRiskProfile[];
}
