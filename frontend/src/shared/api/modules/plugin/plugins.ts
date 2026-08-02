import { ApiResponse, PageData } from "../../types/basic.model";
import apiClient from "../../client";
import { PluginMeta } from "../../types/plugin.model";

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

// export interface ListPluginRequest {
//   page?: number;
//   page_size?: number;
//   author_id?: number;
//   tags?: string[];
//   type?: string;
//   keyword?: string;
//   sort_by?: string;
//   status?: "active" | "inactive" | "all";
// }

export interface PluginListParams {
  enabled?: boolean;
  page?: number;
  page_size?: number;
}

export interface PluginVO {
  id: string;
  name: string;
  authorId: number;
  createdAt: string;
  status: string;
}
export const pluginApi = {
  /** 获取插件列表 */
  list(params?: PluginListParams) {
    return apiClient.get<ApiResponse<PageData<PluginMeta>>>("/plugins", {
      params,
    });
  },

  /** 上传插件文件 */
  upload(file: File, fileType: string = "plugin") {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("file_type", fileType);
    return apiClient.post<ApiResponse<PluginMeta>>("/plugins", formData);
  },
  /** 获取已启用插件（前端运行时加载用） */
  listEnabled() {
    return apiClient.get<ApiResponse<PageData<PluginMeta>>>("/plugins", {
      params: { enabled: true },
    });
  },

  /** 获取单个插件详情 */
  get(id: string) {
    return apiClient.get<ApiResponse<PluginMeta>>(`/plugins/${id}`);
  },

  // /** 创建/安装插件 */
  create(payload: CreatePluginPayload) {
    return apiClient.post<ApiResponse<PluginMeta>>("/plugins", payload);
  },

  /** 更新插件信息 */
  update({ id, ...payload }: UpdatePluginPayload) {
    return apiClient.put<ApiResponse<PluginMeta>>(`/plugins/${id}`, payload);
  },

  /** 启用/禁用插件 */
  toggle(id: string, enabled: boolean) {
    return apiClient.patch<ApiResponse<PluginMeta>>(`/plugins/${id}/toggle`, {
      enabled,
    });
  },

  /** 删除插件 */
  delete(id: string) {
    return apiClient.delete<ApiResponse<null>>(`/plugins/${id}`);
  },
};
