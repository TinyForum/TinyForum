import apiClient from "../client";
import type { PluginMeta } from "@/shared/plugin/types";

export interface CreatePluginPayload {
  name: string;
  version: string;
  description: string;
  author: string;
  scriptUrl: string;
  enabled: boolean;
  slots?: string[];
  config?: Record<string, unknown>;
}

export interface UpdatePluginPayload extends Partial<CreatePluginPayload> {
  id: string;
}

export interface PluginListParams {
  enabled?: boolean;
  page?: number;
  page_size?: number;
}

export const pluginApi = {
  /** 获取插件列表（管理端，全量） */
  list(params?: PluginListParams) {
    return apiClient.get<{ data: PluginMeta[]; total: number }>("/plugins", {
      params,
    });
  },

  /** 获取已启用插件（前端运行时加载用） */
  listEnabled() {
    return apiClient.get<{ data: PluginMeta[] }>("/plugins", {
      params: { enabled: true },
    });
  },

  /** 获取单个插件详情 */
  get(id: string) {
    return apiClient.get<{ data: PluginMeta }>(`/plugins/${id}`);
  },

  /** 创建/安装插件 */
  create(payload: CreatePluginPayload) {
    return apiClient.post<{ data: PluginMeta }>("/plugins", payload);
  },

  /** 更新插件信息 */
  update({ id, ...payload }: UpdatePluginPayload) {
    return apiClient.put<{ data: PluginMeta }>(`/plugins/${id}`, payload);
  },

  /** 启用/禁用插件 */
  toggle(id: string, enabled: boolean) {
    return apiClient.patch<{ data: PluginMeta }>(`/plugins/${id}/toggle`, {
      enabled,
    });
  },

  /** 删除插件 */
  delete(id: string) {
    return apiClient.delete(`/plugins/${id}`);
  },
};
