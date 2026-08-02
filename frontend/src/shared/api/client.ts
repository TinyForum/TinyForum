// lib/api/client.ts
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from "axios";
import { routing } from "@/i18n/routing";

const API_BASE_URL = "/api/v1";

function resolveLoginPath(): string {
  const segments = window.location.pathname.split("/");
  const first = segments[1] ?? "";
  // 已带 locale 前缀（如 /zh-CN/auth/login）时直接复用
  if (routing.locales.includes(first as (typeof routing.locales)[number])) {
    return `/${first}/auth/login`;
  }
  // 无前缀（如 /auth/login 或 /），用默认 locale 补齐
  return `/${routing.defaultLocale}/auth/login`;
}

// const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;
function createClient(config?: AxiosRequestConfig): AxiosInstance {
  const instance = axios.create({
    baseURL: API_BASE_URL,
    headers: { "Content-Type": "application/json" },
    timeout: 10_000,
    withCredentials: true, // Cookie 自动携带
    ...config,
  });

  // 请求拦截：FormData 交由浏览器生成 boundary，移除手动 Content-Type
  instance.interceptors.request.use((config) => {
    if (config.data instanceof FormData) {
      delete config.headers?.["Content-Type"];
    }
    return config;
  });

  // 响应拦截：401 处理
  instance.interceptors.response.use(
    (res) => res,
    (err: AxiosError) => {
      const url = err.config?.url ?? "";

      // ✅ 登出接口：直接返回成功，不触发重定向
      if (url.includes("/auth/logout")) {
        return Promise.resolve({ data: { code: 0 } });
      }

      // ✅ 会话检查/账户状态接口：仅拒绝，由调用方自行处理
      //    （auth store refreshUser / PostLoginHandler / useUserRole 等）
      if (
        url.includes("/users/me/role") ||
        url.includes("/auth/account/") ||
        url.includes("/auth/me")
      ) {
        return Promise.reject(err);
      }

      if (err.response?.status === 401) {
        window.location.href = resolveLoginPath();
      }
      return Promise.reject(err);
    },
  );

  return instance;
}

const apiClient = createClient();
export default apiClient;
export { createClient, API_BASE_URL };
