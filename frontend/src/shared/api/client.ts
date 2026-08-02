// lib/api/client.ts
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from "axios";
import { routing } from "@/i18n/routing";
import { useAuthStore } from "@/store/auth";

const API_BASE_URL = "/api/v1";

// 解析带 locale 的登录路径
function resolveLoginPath(): string {
  const segments = window.location.pathname.split("/");
  const first = segments[1] ?? "";
  if (routing.locales.includes(first as (typeof routing.locales)[number])) {
    return `/${first}/auth/login`;
  }
  return `/${routing.defaultLocale}/auth/login`;
}

// 防并发锁
let isHandling401 = false;

function createClient(config?: AxiosRequestConfig): AxiosInstance {
  const instance = axios.create({
    baseURL: API_BASE_URL,
    headers: { "Content-Type": "application/json" },
    timeout: 10_000,
    withCredentials: true,
    ...config,
  });

  // 请求拦截：FormData 处理
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
      // 跳过登出接口本身的 401（避免循环）
      if (err.config?.url?.includes("/auth/logout")) {
        return Promise.resolve({ data: { code: 0 } });
      }

      if (err.response?.status === 401) {
        // 如果正在处理 401，直接拒绝，避免并发重复操作
        if (isHandling401) {
          return Promise.reject(err);
        }

        const currentPath = window.location.pathname;
        // 如果已经在登录页或注册页，不处理（防止无限重定向）
        if (
          currentPath.includes("/auth/login") ||
          currentPath.includes("/auth/register")
        ) {
          return Promise.reject(err);
        }

        // 加锁
        isHandling401 = true;

        // 同步清除认证状态（清除 localStorage、sessionStorage、cookie）
        useAuthStore.getState().logout();

        // 重定向到登录页（replace 避免回退到当前页）
        const loginUrl = resolveLoginPath();
        window.location.replace(loginUrl);

        // 解锁（但页面已跳转，这里可能不会执行，但为了安全，定时解锁）
        setTimeout(() => {
          isHandling401 = false;
        }, 500);

        return Promise.reject(err);
      }

      return Promise.reject(err);
    },
  );

  return instance;
}

const apiClient = createClient();
export default apiClient;
export { createClient, API_BASE_URL };
