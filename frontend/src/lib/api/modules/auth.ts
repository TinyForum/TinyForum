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

  login: (data: LoginPayload) =>
    apiClient.post<ApiResponse<AuthResult>>("/auth/login", data),

  me: () =>
    apiClient.get<ApiResponse<User>>("/auth/me"),
};