/**
 * api/modules/admin.ts
 * 管理员 API（用户、帖子、板块管理）
 */

import apiClient from "../../client";
import { ApiResponse, PageData } from "../../types/basic.model";
import { UserDO } from "../../types/user.model";

// ========== 类型定义 ==========
export interface ResetPasswordResponse {
  message: string;
  operator_id: number;
  user_id: number;
}

// ========== 管理员 API ==========
export const adminUsersApi = {
  // ── 用户管理 ──────────────────────────────────────────────────────────────
  /** 获取用户列表（分页，可选关键词搜索） */
  listUsers: (params?: {
    page?: number;
    page_size?: number;
    keyword?: string;
  }) =>
    apiClient.get<ApiResponse<PageData<UserDO>>>("/admin/users", { params }),

  /** 设置用户激活状态 */
  setUserActive: (id: number, isActive: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/active`, {
      is_active: isActive,
    }),

  /** 设置用户封禁状态 */
  setUserBlocked: (id: number, isBlocked: boolean) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/blocked`, {
      is_blocked: isBlocked,
    }),

  /** 设置用户角色 */
  setUserRole: (id: number, role: string) =>
    apiClient.put<ApiResponse<null>>(`/admin/users/${id}/role`, { role }),

  /** 删除用户（软删除） */
  deleteUser: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/admin/users/${id}`),

  /** 重置用户密码（返回新密码或操作结果） */
  resetUserPassword: (id: number) =>
    apiClient.post<ApiResponse<ResetPasswordResponse>>(
      `/admin/users/${id}/reset-password`,
    ),
};
