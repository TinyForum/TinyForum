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

  update: (file: string, data: Record<string, unknown>) =>
    apiClient.put<ApiResponse<null>>(`/admin/config/${file}`, data, {
      withCredentials: true,
    }),

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
