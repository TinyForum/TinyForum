import { apiClient } from "@/shared/api";

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
  uploadCommentFile(commentId: string | number, file: File) {
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

  /** 获取当前用户文件列表 */
  getUserFiles(params?: { page?: number; page_size?: number }) {
    return apiClient.get("/attachments/users/me/files", { params });
  },

  /** 获取帖子文件信息 */
  getPostFile(fileId: string) {
    return apiClient.get(`/attachments/post/${fileId}`);
  },

  /** 获取评论文件信息 */
  getCommentFile(fileId: string) {
    return apiClient.get(`/attachments/comment/${fileId}`);
  },

  /** 获取插件文件信息 */
  getPluginFile(fileId: string) {
    return apiClient.get(`/attachments/plugin/${fileId}`);
  },

  /** 删除帖子文件 */
  deletePostFile(fileId: string) {
    return apiClient.delete(`/attachments/post/${fileId}`);
  },

  /** 删除评论文件 */
  deleteCommentFile(fileId: string) {
    return apiClient.delete(`/attachments/comment/${fileId}`);
  },

  /** 删除插件文件 */
  deletePluginFile(fileId: string) {
    return apiClient.delete(`/attachments/plugin/${fileId}`);
  },

  /** 公开访问文件（无需认证） */
  serveFile(fileId: string) {
    return apiClient.get(`/files/${fileId}`, { responseType: "blob" });
  },
};
