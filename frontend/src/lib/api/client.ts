/**
 * api/client.ts
 * Axios 基础客户端 —— 统一拦截器、错误处理、Token 管理
 */

import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from "axios";

// ─── 常量 ─────────────────────────────────────────────────────────────────────

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

const TOKEN_KEY = "bbs_token";
const USER_KEY = "bbs_user";

// ─── 工厂函数 ─────────────────────────────────────────────────────────────────

function createClient(config?: AxiosRequestConfig): AxiosInstance {
  const instance = axios.create({
    baseURL: API_BASE_URL,
    headers: { "Content-Type": "application/json" },
    timeout: 10_000,
    withCredentials: true,
    ...config,
  });

  // ── 请求拦截：附加 Bearer Token ──────────────────────────────────────────────
  instance.interceptors.request.use((req) => {
    if (typeof window !== "undefined") {
      const token = localStorage.getItem(TOKEN_KEY);
      if (token) req.headers.Authorization = `Bearer ${token}`;
    }
    return req;
  });

  // ── 响应拦截：统一 401 处理 ───────────────────────────────────────────────────
  instance.interceptors.response.use(
    (res) => res,
    (err: AxiosError) => {
      if (err.response?.status === 401 && typeof window !== "undefined") {
        localStorage.removeItem(TOKEN_KEY);
        localStorage.removeItem(USER_KEY);
        window.location.href = "/auth/login";
      }
      return Promise.reject(err);
    }
  );

  return instance;
}

// ─── 导出单例 ─────────────────────────────────────────────────────────────────

/** 默认客户端实例（整个应用共用） */
const apiClient = createClient();
export default apiClient;
export { createClient, API_BASE_URL, TOKEN_KEY, USER_KEY };