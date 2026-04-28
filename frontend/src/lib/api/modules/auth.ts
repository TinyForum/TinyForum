/**
 * api/modules/auth.ts
 */

import { ApiResponse, AuthResult, User } from "@/shared/api/types";
import apiClient from "../client";
import ForgotPasswordPage from "@/app/[locale]/auth/forgot-password/page";

export interface RegisterPayload {
  username: string;
  email: string;
  password: string;
}

export interface LoginPayload {
  email: string;
  password: string;
}

export const authApi = {
  register: (data: RegisterPayload) =>
    apiClient.post<ApiResponse<AuthResult>>("/auth/register", data),

  // 登录：后端通过 Set-Cookie 设置 HttpOnly Cookie
  login: (data: LoginPayload) =>
    apiClient.post<ApiResponse<AuthResult>>("/auth/login", data, {
      withCredentials: true,
    }),

  // 获取当前用户：Cookie 会自动携带
  me: () =>
    apiClient.get<ApiResponse<User>>("/auth/me", {
      withCredentials: true,
    }),

  // 登出：清除 Cookie
  logout: () =>
    apiClient.post<ApiResponse<null>>("/auth/logout", null, {
      withCredentials: true,
    }),

  // 请求注销账户（软删除）
  deleteAccount: (data?: { confirm: string; password?: string }) =>
    apiClient.delete<ApiResponse<null>>("/auth/account", {
      withCredentials: true,
      data,
    }),

  // 取消注销（恢复账户）
  cancelDeletion: () =>
    apiClient.post<ApiResponse<null>>("/auth/account/restore", null, {
      withCredentials: true,
    }),

  // 确认永久删除（硬删除）
  confirmDeletion: (data?: { confirm: string; password?: string }) =>
    apiClient.delete<ApiResponse<null>>("/auth/account/permanent", {
      withCredentials: true,
      data,
    }),
  // 获取注销状态
  getDeletionStatus: () =>
    apiClient.get<
      ApiResponse<{
        is_deleted: boolean;
        deleted_at?: string;
        can_restore: boolean;
        remaining_days?: number;
      }>
    >("/auth/account/deletion", {
      withCredentials: true,
    }),
  // 忘记密码
  forgotPassword: (data: { email: string }) =>
    apiClient.post<ApiResponse<null>>("/auth/password/forgot", data),
  // 重置密码
  resetPassword: (data: { token: string; password: string }) =>
    apiClient.post<ApiResponse<null>>("/auth/password/reset", data),
};
