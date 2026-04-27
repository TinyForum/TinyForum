// lib/api/client.ts
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from "axios";

// const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "/api/v1";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";
function createClient(config?: AxiosRequestConfig): AxiosInstance {
  console.log("API_BASE_URL: ", API_BASE_URL, "  config: ", config);
  const instance = axios.create({
    baseURL: API_BASE_URL,
    headers: { "Content-Type": "application/json" },
    timeout: 10_000,
    withCredentials: true, // Cookie 自动携带
    ...config,
  });

  // 响应拦截：401 处理
  instance.interceptors.response.use(
    (res) => res,
    (err: AxiosError) => {
      if (err.response?.status === 401 && typeof window !== "undefined") {
        const pathname = window.location.pathname;
        const isLoginPage = pathname.includes("/auth/login");
        if (!isLoginPage) {
          // 提取当前 locale 前缀，保持跳转一致
          const localeMatch = pathname.match(/^\/(zh-CN|en|ja)/);
          const locale = localeMatch ? localeMatch[1] : "zh-CN";
          window.location.href = `/${locale}/auth/login`;
        }
      }
      return Promise.reject(err);
    },
  );

  return instance;
}

const apiClient = createClient();
export default apiClient;
export { createClient, API_BASE_URL };
