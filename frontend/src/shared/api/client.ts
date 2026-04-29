// lib/api/client.ts
import axios, { AxiosError, AxiosInstance, AxiosRequestConfig } from "axios";

const API_BASE_URL = "/api/v1";

// const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL;
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
      // ✅ 跳过登出接口
      if (err.config?.url?.includes("/auth/logout")) {
        return Promise.resolve({ data: { code: 0 } });
      }

      if (err.response?.status === 401) {
        window.location.href = `/auth/login`;
      }
      return Promise.reject(err);
    },
  );

  return instance;
}

const apiClient = createClient();
export default apiClient;
export { createClient, API_BASE_URL };
