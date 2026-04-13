/**
 * api/modules/users.ts
 */

import apiClient from "../client";
import type { ApiResponse, PageData, User } from "../types";

export interface UpdateProfilePayload {
  bio?: string;
  avatar?: string;
}

export const userApi = {
  getProfile: (id: number) =>
    apiClient.get<ApiResponse<User>>(`/users/${id}`),

  updateProfile: (data: UpdateProfilePayload) =>
    apiClient.put<ApiResponse<User>>("/users/profile", data),

  follow: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/users/${id}/follow`),

  unfollow: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/users/${id}/follow`),

  leaderboard: (limit?: number) =>
    apiClient.get<ApiResponse<User[]>>("/users/leaderboard", {
      params: { limit },
    }),

  follwowers: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<User>>>(`/users/${id}/followers`, { params }),
  following: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<User>>>(`/users/${id}/following`, { params }),
  // ── Admin ─────────────────────────────────────────────────────────────────────
  adminList: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    apiClient.get<ApiResponse<PageData<User>>>("/admin/users", { params }),

  adminSetActive: (id: number, active: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/active`, { active }),

  adminSetBlocked: (id: number, blocked: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/blocked`, { blocked }),

  adminSetRole: (id: number, role: string) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/role`, { role }),
};