/**
 * api/modules/users.ts
 */

import { UserRoleType } from "@/type/roles.types";
import apiClient from "../client";
import type { ApiResponse, PageData, User } from "../types";

export interface UpdateProfilePayload {
  bio?: string;
  avatar?: string;
}

export interface RoleResponse {
  user_id: number;
  role: UserRoleType;
}

export const userApi = {
  // 获取用户信息
  getProfile: (id: number) => apiClient.get<ApiResponse<User>>(`/users/${id}`),
  // 获取当前用户的角色
  getMeRole: () => apiClient.get<ApiResponse<RoleResponse>>("/users/me/role"),

  // 更新用户信息
  updateProfile: (data: UpdateProfilePayload) =>
    apiClient.put<ApiResponse<User>>("/users/profile", data),

  // 更改密码
  changePassword: (params: { old_password: string; new_password: string }) =>
    apiClient.patch<ApiResponse<null>>("/users/password", params),
  // 获取用户关注
  follow: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/users/${id}/follow`),

  unfollow: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/users/${id}/follow`),

  leaderboard: (limit?: number) =>
    apiClient.get<ApiResponse<User[]>>("/users/leaderboard", {
      params: { limit },
    }),

  follwowers: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<User>>>(`/users/${id}/followers`, {
      params,
    }),
  following: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<User>>>(`/users/${id}/following`, {
      params,
    }),
  // ── Admin ─────────────────────────────────────────────────────────────────────
  adminList: (params?: {
    page?: number;
    page_size?: number;
    keyword?: string;
  }) => apiClient.get<ApiResponse<PageData<User>>>("/admin/users", { params }),

  adminSetActive: (id: number, active: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/active`, { active }),

  adminSetBlocked: (id: number, blocked: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/blocked`, { blocked }),

  adminSetRole: (id: number, role: string) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/role`, { role }),
};
