// lib/api/client.ts
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from "axios";

// const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "/api/v1";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
function createClient(config?: AxiosRequestConfig): AxiosInstance {
  const instance = axios.create({
    baseURL: API_BASE_URL,
    headers: { "Content-Type": "application/json" },
    timeout: 10_000,
    withCredentials: true, // Cookie 自动携带，不需要手动附加 token
    ...config,
  });

  // 请求拦截：不再手动附加 token，Cookie 由浏览器自动处理
  // （如果将来需要 CSRF token 可以在这里加）

  // 响应拦截：401 处理
  instance.interceptors.response.use(
    (res) => res,
    (err: AxiosError) => {
      if (err.response?.status === 401 && typeof window !== "undefined") {
        // 不要直接跳转，避免在登录页也触发重定向循环
        const isLoginPage = window.location.pathname.includes('/auth/login');
        if (!isLoginPage) {
          window.location.href = "/auth/login";
        }
      }
      return Promise.reject(err);
    }
  );

  return instance;
}

const apiClient = createClient();
export default apiClient;
export { createClient, API_BASE_URL };