/**
 * api/modules/users.ts
 * 用户相关 API（不含管理员功能）
 */

import apiClient from "../client";
import { UserRoleType } from "@/shared/type/roles.types";
import { UserDO } from "../types/user.model";
import { ApiResponse, PageData } from "../types/basic.model";

// ========== 请求/响应类型 ==========
export interface UpdateProfilePayload {
  bio?: string;
  avatar?: string;
}

export interface RoleResponse {
  user_id: number;
  role: UserRoleType;
}

export interface LeaderboardRequest {
  limit?: number;
}

export interface LeaderboardItemResponse {
  id: number;
  username: string;
  avatar: string;
  score: number;
  rank: number;
  bio: string;
}

// ========== API 方法 ==========
export const userApi = {
  // ---------- 用户公开信息 ----------
  /** 获取指定用户公开资料 (无需登录) */
  getProfile: (id: number) =>
    apiClient.get<ApiResponse<UserDO>>(`/users/${id}`),

  /** 获取当前用户的角色 (需登录) */
  getMeRole: () => apiClient.get<ApiResponse<RoleResponse>>("/users/me/role"),

  /** 更新当前用户资料 (需登录) */
  updateProfile: (data: UpdateProfilePayload) =>
    apiClient.put<ApiResponse<UserDO>>("/users/mine/profile", data),

  // ---------- 关注/取关 ----------
  /** 关注用户 (需登录) */
  follow: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/users/${id}/follow`),

  /** 取消关注 (需登录) */
  unfollow: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/users/${id}/follow`),

  // ---------- 榜单 ----------
  /** 简洁排行榜 (积分排行) */
  getLeaderboardSimple: (params?: LeaderboardRequest) =>
    apiClient.get<ApiResponse<LeaderboardItemResponse[]>>(
      "/users/leaderboard/simple",
      { params },
    ),

  /** 详细排行榜 (含更多字段) */
  getLeaderboardDetail: (params?: LeaderboardRequest) =>
    apiClient.get<ApiResponse<LeaderboardItemResponse[]>>(
      "/users/leaderboard/detail",
      { params },
    ),

  // ---------- 粉丝/关注列表 ----------
  /** 获取用户的粉丝列表 (分页) */
  getFollowers: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<UserDO>>>(`/users/${id}/followers`, {
      params,
    }),

  /** 获取用户关注的列表 (分页) */
  getFollowing: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<UserDO>>>(`/users/${id}/following`, {
      params,
    }),
};
