import apiClient from "../client";
import { ApiResponse } from "../types/basic.model";

export const uploadApi = {
  // ========== 上传相关 API ==========
  /** 上传帖子文件 */
  uploadPostFile(postId: string | number, file: File) {
    const formData = new FormData();
    formData.append("file", file);
    return apiClient.post<{ data: string }>(
      `/attachments/post/${postId}`,
      formData,
      {
        headers: { "Content-Type": "multipart/form-data" },
      },
    );
  },

  /** 上传评论文件 */
  uploadCommentFile: (commentId: string | number, file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    return apiClient.post<{ data: string }>(
      `/attachments/comment/${commentId}`,
      formData,
      {
        headers: { "Content-Type": "multipart/form-data" },
      },
    );
  },

  /** 上传插件文件 */
  uploadPluginFile(file: File, fileType: string = "plugin") {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("file_type", fileType);
    return apiClient.post<{ data: string }>("/attachments/plugin", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    });
  },

  /** 用户上传头像 */
  uploadAvatar: (file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("type", "avatar");
    return apiClient.post<ApiResponse<UploadResponse>>(
      "/attachments",
      formData,
      {
        headers: { "Content-Type": "multipart/form-data" },
      },
    );
  },

  /** 获取当前用户文件列表 */
  getUserFiles: (params?: { page?: number; page_size?: number }) =>
    apiClient.get("/attachments/users/me/files", { params }),

  /** 获取帖子文件信息 */
  getPostFile: (fileId: string) => apiClient.get(`/attachments/post/${fileId}`),

  /** 获取评论文件信息 */
  getCommentFile: (fileId: string) =>
    apiClient.get(`/attachments/comment/${fileId}`),

  /** 获取插件文件信息 */
  getPluginFile: (fileId: string) =>
    apiClient.get(`/attachments/plugin/${fileId}`),

  // 获取用户上传的插件信息
  getUserPlugins: (params?: { page?: number; page_size?: number }) =>
    apiClient.get("/attachments/plugin/users/me", { params }),

  /** 删除帖子文件 */
  deletePostFile: (fileId: string) =>
    apiClient.delete(`/attachments/post/${fileId}`),

  /** 删除评论文件 */
  deleteCommentFile: (fileId: string) =>
    apiClient.delete(`/attachments/comment/${fileId}`),

  /** 删除插件文件 */
  deletePluginFile: (fileId: string) =>
    apiClient.delete(`/attachments/plugin/${fileId}`),

  /** 公开访问文件（无需认证） */
  serveFile: (fileId: string) =>
    apiClient.get(`/files/${fileId}`, { responseType: "blob" }),
};

export interface UploadResponse {
  file_id: string; // 存储标识
  url: string; // 访问URL
  original_name: string;
  size: number;
  mime_type: string;
}
