/**
 * api/modules/auth.ts
 */

import { ApiResponse, AuthResult, User } from "@/shared/api/types";
import apiClient from "../client";

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
    apiClient.post<ApiResponse<AuthResult>>("/auth/login", data),

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
  forgotPassword: (data: ForgotPasswordRequest) =>
    apiClient.post<ApiResponse<null>>("/auth/password/forgot", data),
  // 重置密码
  resetPassword: (data: ResetPasswordRequest) =>
    apiClient.put<ApiResponse<ResetPasswordResponse>>(
      "/auth/password/reset",
      data,
    ),
  // 通过 token 重置密码
  resetPasswordWithToken: (data: ResetPasswordWithTokenRequest) =>
    apiClient.put<ApiResponse<ResetPasswordWithTokenResponse>>(
      "/auth/password/reset-withtoken",
      data,
    ),
  // 验证token
  validateToken: (data: ValidateTokenRequest) =>
    apiClient.get<ApiResponse<ValidateTokenResponse>>(
      `/auth/password/validate-token?token=${data.token}`,
    ),
  // 更改密码
  changePassword: (params: { old_password: string; new_password: string }) =>
    apiClient.put<ApiResponse<null>>("/auth/account/password", params),
};

// 请求
export interface ForgotPasswordRequest {
  email: string;
}

// 重置密码
export interface ResetPasswordRequest {
  password: string;
}
// 通过 token 重置密码
export interface ResetPasswordWithTokenRequest {
  password: string;
  token: string;
}
export interface ResetPasswordWithTokenResponse {
  success: boolean;
}
export interface ResetPasswordResponse {
  success: boolean;
}

// 验证 token
export interface ValidateTokenRequest {
  token: string;
}

export interface ValidateTokenResponse {
  valid: boolean;
}
