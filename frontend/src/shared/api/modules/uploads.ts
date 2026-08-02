import apiClient from "../client";
import { ApiResponse } from "../types/basic.model";

// ========== 类型定义 ==========

/** 上传文件响应（POST /attachments） */
export interface UploadResponse {
  file_id: string; // 存储标识
  url: string; // 访问URL
  original_name: string;
  size: number;
  mime_type: string;
}

/** 文件元信息（GET /attachments/{file_id}） */
export interface FileInfoResponse {
  file_id: string;
  original_name: string;
  size: number;
  mime_type: string;
  url: string;
  file_type: string;
  created_at: string;
}

/** 用户文件列表（GET /attachments/user/me） */
export interface FileListResponse {
  list: FileInfoResponse[];
  total: number;
}

// ========== API 方法 ==========
export const uploadApi = {
  /**
   * 通用上传（type 为业务类型：post_image / comment_file / avatar / plugin / post_cover / topic_cover）
   * @param params.type - 业务类型，决定后续关联字段是否必需
   * @param params.post_id - 当 type = 'post_image' 时必需
   * @param params.reply_id - 当 type = 'comment_file' 时必需
   */
  uploadFile(params: {
    file: File;
    type: string;
    post_id?: string | number;
    reply_id?: string | number;
  }) {
    const formData = new FormData();
    formData.append("file", params.file);
    formData.append("type", params.type);
    if (params.post_id !== undefined) {
      formData.append("post_id", String(params.post_id));
    }
    if (params.reply_id !== undefined) {
      formData.append("reply_id", String(params.reply_id));
    }
    return apiClient.post<ApiResponse<UploadResponse>>(
      "/attachments",
      formData,
    );
  },

  /** 上传帖子图片（关联 post_id） */
  uploadPostFile(postId: string | number, file: File) {
    return this.uploadFile({ file, type: "post_image", post_id: postId });
  },

  /** 上传评论附件（关联 reply_id） */
  uploadCommentFile(commentId: string | number, file: File) {
    return this.uploadFile({ file, type: "comment_file", reply_id: commentId });
  },

  /** 上传封面/插件等通用文件（type 如 post_cover / topic_cover / plugin） */
  uploadPluginFile(file: File, fileType: string = "plugin") {
    return this.uploadFile({ file, type: fileType });
  },

  /** 上传头像 */
  uploadAvatar(file: File) {
    return this.uploadFile({ file, type: "avatar" });
  },

  /** 获取当前用户文件列表（分页，可按 file_type 过滤） */
  getUserFiles(params?: {
    page?: number;
    page_size?: number;
    file_type?: string;
  }) {
    return apiClient.get<ApiResponse<FileListResponse>>("/attachments/user/me", {
      params,
    });
  },

  /** 获取当前用户的插件文件列表（file_type=plugin） */
  getUserPlugins(params?: { page?: number; page_size?: number }) {
    return apiClient.get<ApiResponse<FileListResponse>>("/attachments/user/me", {
      params: { ...params, file_type: "plugin" },
    });
  },

  /** 获取帖子文件信息 */
  getPostFile(fileId: string) {
    return apiClient.get<ApiResponse<FileInfoResponse>>(
      `/attachments/${fileId}`,
    );
  },

  /** 获取评论文件信息 */
  getCommentFile(fileId: string) {
    return apiClient.get<ApiResponse<FileInfoResponse>>(
      `/attachments/${fileId}`,
    );
  },

  /** 获取插件文件信息 */
  getPluginFile(fileId: string) {
    return apiClient.get<ApiResponse<FileInfoResponse>>(
      `/attachments/${fileId}`,
    );
  },

  /** 删除帖子文件（仅所有者） */
  deletePostFile(fileId: string) {
    return apiClient.delete<ApiResponse<{ message: string }>>(
      `/attachments/${fileId}`,
    );
  },

  /** 删除评论文件（仅所有者） */
  deleteCommentFile(fileId: string) {
    return apiClient.delete<ApiResponse<{ message: string }>>(
      `/attachments/${fileId}`,
    );
  },

  /** 删除插件文件（仅所有者） */
  deletePluginFile(fileId: string) {
    return apiClient.delete<ApiResponse<{ message: string }>>(
      `/attachments/${fileId}`,
    );
  },

  /** 公开访问文件（无需认证） */
  serveFile(fileId: string) {
    return apiClient.get<Blob>(`/files/${fileId}`, { responseType: "blob" });
  },
};
