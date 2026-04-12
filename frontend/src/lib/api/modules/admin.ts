/**
 * api/modules/admin.ts
 */

import apiClient from "../client";
import type { ApiResponse, PageData, User, Post, Board } from "../types";

export const adminApi = {
  // ── 用户管理 ──────────────────────────────────────────────────────────────────
  listUsers: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    apiClient.get<ApiResponse<PageData<User>>>("/admin/users", { params }),

  setUserActive: (id: number, active: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/active`, { active }),

  setUserBlocked: (id: number, blocked: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/blocked`, { blocked }),

  setUserRole: (id: number, role: string) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/role`, { role }),

  // ── 帖子管理 ──────────────────────────────────────────────────────────────────
  listPosts: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    apiClient.get<ApiResponse<PageData<Post>>>("/admin/posts", { params }),

  togglePin: (id: number) =>
    apiClient.put<ApiResponse<null>>(`/admin/posts/${id}/pin`),

  // ── 板块管理 ──────────────────────────────────────────────────────────────────
  listBoards: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Board>>>("/admin/boards", { params }),
};