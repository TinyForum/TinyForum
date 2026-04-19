/**
 * api/modules/admin.ts
 */

import apiClient from "../client";
import type { ApiResponse, PageData, User, Post, Board } from "../types";

// MARK: 请求/响应体定义
// 重置密码
export interface AdminResetPasswordRequest {
id: number ;
}

export interface ResetPasswordResponse {
  message: string;
    operator_id:number;
    user_id: number
}

// MARK: API 定义
export const adminApi = {
  // ── 用户管理 ──────────────────────────────────────────────────────────────────
  /**
   *  列出所有用户
   * @param params 
   * @returns 
   */
  listUsers: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    apiClient.get<ApiResponse<PageData<User>>>("/admin/users", { params }),

  /** 设置用户激活状态（ 主动状态）
   * @param id 用户 id
   * @param is_active bool
   * @returns
   */ 
  setUserActive: (id: number, is_active: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/active`, { is_active }),

  /**
   * 设置用户封禁状态
   * @param id 用户 id
   * @param blocked 
   * @returns 
   */
  setUserBlocked: (id: number, is_blocked: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/blocked`, { is_blocked }),

  // 设置用户角色
  setUserRole: (id: number, role: string) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/role`, { role }),
  // 删除用户(软删除)
  deleteUser: (id: number) => apiClient.delete<ApiResponse<null>>(`/admin/users/${id}`),
  // 重置用户密码
  resetUserPassword: (id: number) => apiClient.post<ApiResponse<ResetPasswordResponse>>(`/admin/users/${id}/reset-password`),

  // ── 帖子管理 ──────────────────────────────────────────────────────────────────
  // 列出所有帖子
  listPosts: (params?: { page?: number; page_size?: number; keyword?: string }) =>
    apiClient.get<ApiResponse<PageData<Post>>>("/admin/posts", { params }),

  // 置顶
  togglePin: (id: number) =>
    apiClient.put<ApiResponse<null>>(`/admin/posts/${id}/pin`),
  // 获取所有待审核帖子
  listPendingPosts: (params?: { page?: number; page_size?: number ,keyword?: string}) => apiClient.get<ApiResponse<PageData<Post>>>("/admin/posts/pending", { params }),
  // 标记审核状态
  reviewPosts: (id: number, data: { status: string }) =>
  apiClient.put<ApiResponse<null>>(`/admin/posts/${id}/review`, data),
  
  // ── 板块管理 ──────────────────────────────────────────────────────────────────
  listBoards: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Board>>>("/admin/boards", { params }),
};