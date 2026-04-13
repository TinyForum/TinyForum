/**
 * api/modules/auth.ts
 */

import apiClient from "../client";
import type { ApiResponse, AuthResult, User } from "../types";

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
      withCredentials: true, // 重要：允许接收和发送 Cookie
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
};