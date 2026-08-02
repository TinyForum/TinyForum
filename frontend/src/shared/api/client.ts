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
      // ✅ 跳过登出接口
      if (err.config?.url?.includes("/auth/logout")) {
        return Promise.resolve({ data: { code: 0 } });
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
