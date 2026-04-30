/**
 * api/modules/users.ts
 */

import { UserRoleType } from "@/shared/type/roles.types";
import apiClient from "../client";
import { ApiResponse, PageData, User } from "@/shared/api/types";

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
  // fields?: string;
}

export interface LeaderboardItemResponse {
  id: number;
  username: string;
  avatar: string;
  score: number;
  rank: number;
  bio: string;
}

/**
 * @deprecated 使用新的 API: userApi / adminusersApi 此 api 存在调用安全问题
 */
export const userAPI = {
  // 获取用户信息
  getProfile: (id: number) => apiClient.get<ApiResponse<User>>(`/users/${id}`),
  // 获取当前用户的角色
  getMeRole: () => apiClient.get<ApiResponse<RoleResponse>>("/users/me/role"),

  // 更新用户信息
  updateProfile: (data: UpdateProfilePayload) =>
    apiClient.put<ApiResponse<User>>("/users/profile", data),

  // 获取用户关注
  follow: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/users/${id}/follow`),

  unfollow: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/users/${id}/follow`),

  // 积分排行精简
  getLeaderboardSimple: (params?: LeaderboardRequest) =>
    apiClient.get<ApiResponse<LeaderboardItemResponse[]>>(
      "/users/leaderboard/simple",
      {
        params,
      },
    ),
  // 积分排行详情
  getLeaderboardDetail: (params?: LeaderboardRequest) =>
    apiClient.get<ApiResponse<LeaderboardItemResponse[]>>(
      "/users/leaderboard/detail",
      {
        params,
      },
    ),

  follwowers: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<User>>>(`/users/${id}/followers`, {
      params,
    }),
  following: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<User>>>(`/users/${id}/following`, {
      params,
    }),
  // ── Admin ─────────────────────────────────────────────────────────────────────
};
