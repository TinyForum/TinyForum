import apiClient from "../client";
import { ApiResponse } from "../types/basic.model";

export interface ConfigFileItem {
  name: string;
  path: string;
  size: number;
  updatedAt: string;
}

export interface ConfigReloadResult {
  success: boolean;
  message: string;
  reloadedFiles: string[];
  errors: string[];
}

export interface ConfigKVResponse {
  name: string;
  format: "kv";
  kv: Record<string, unknown>;
}

export const configApi = {
  list: () =>
    apiClient.get<ApiResponse<ConfigFileItem[]>>("/admin/config/list", {
      withCredentials: true,
    }),

  get: (file: string) =>
    apiClient.get<ApiResponse<Record<string, unknown>>>(
      `/admin/config/${file}`,
      { withCredentials: true },
    ),

  /** 获取配置文件 KV 格式 */
  getKV: (file: string) =>
    apiClient.get<ApiResponse<ConfigKVResponse>>(
      `/admin/config/${file}/kv`,
      { withCredentials: true },
    ),

  update: (file: string, data: Record<string, unknown>) =>
    apiClient.put<ApiResponse<null>>(`/admin/config/${file}`, data, {
      withCredentials: true,
    }),

  /** 通过 KV 键值对更新配置 */
  updateKV: (file: string, config: Record<string, string>) =>
    apiClient.put<ApiResponse<{ message: string }>>(
      `/admin/config/${file}/kv`,
      { config },
      { withCredentials: true },
    ),

  reload: () =>
    apiClient.post<ApiResponse<ConfigReloadResult>>(
      "/admin/config/reload",
      null,
      { withCredentials: true },
    ),

  getHistory: () =>
    apiClient.get<
      ApiResponse<
        { file: string; updatedAt: string; operator: string }[]
      >
    >("/admin/config/history", { withCredentials: true }),
};
