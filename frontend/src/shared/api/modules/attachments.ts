import apiClient from "../client";
import { ApiResponse } from "../types/basic.model";

// ========== 类型定义 ==========
/** 上传文件请求参数（业务类型） */
export type UploadFileType =
  | "post_image"
  | "comment_file"
  | "avatar"
  | "plugin"
  | "post_cover"
  | "topic_cover";

/** 上传文件响应数据 */
export interface UploadFileResponse {
  file_id: string; // 文件唯一标识
  url: string; // 访问地址
  original_name: string; // 原始文件名
  size: number; // 文件大小（字节）
  mime_type: string; // MIME 类型
}

/** 文件元信息（GET /attachments/{file_id} 返回） */
export interface FileInfoResponse {
  file_id: string;
  original_name: string;
  size: number;
  mime_type: string;
  url: string;
  file_type: string; // 业务类型
  created_at: string;
}

/** 用户文件列表响应 */
export interface UserFileListResponse {
  list: FileInfoResponse[];
  total: number;
}

// ========== API 方法 ==========
export const attachmentApi = {
  /**
   * 上传文件（通用）
   * @param params - 表单参数
   * @param params.file - 文件对象
   * @param params.type - 业务类型，决定后续关联字段是否必需
   * @param params.post_id - 当 type = 'post_image' 时必需
   * @param params.reply_id - 当 type = 'comment_file' 时必需
   */
  uploadFile(params: {
    file: File;
    type: UploadFileType;
    post_id?: number;
    reply_id?: number;
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
    return apiClient.post<ApiResponse<UploadFileResponse>>(
      "/attachments",
      formData,
    );
  },

  /**
   * 便捷方法：上传帖子图片
   */
  uploadPostImage(file: File, postId: number) {
    return this.uploadFile({ file, type: "post_image", post_id: postId });
  },

  /**
   * 便捷方法：上传评论附件
   */
  uploadCommentFile(file: File, replyId: number) {
    return this.uploadFile({ file, type: "comment_file", reply_id: replyId });
  },

  /**
   * 便捷方法：上传头像
   */
  uploadAvatar(file: File) {
    return this.uploadFile({ file, type: "avatar" });
  },

  /**
   * 便捷方法：上传插件
   */
  uploadPlugin(file: File) {
    return this.uploadFile({ file, type: "plugin" });
  },

  /**
   * 获取文件元信息
   * @param fileId - 文件ID
   */
  getFile(fileId: string) {
    return apiClient.get<ApiResponse<FileInfoResponse>>(
      `/attachments/${fileId}`,
    );
  },

  /**
   * 删除文件（仅文件所有者）
   * @param fileId - 文件ID
   */
  deleteFile(fileId: string) {
    return apiClient.delete<ApiResponse<{ message: string }>>(
      `/attachments/${fileId}`,
    );
  },

  /**
   * 公开访问文件（返回文件二进制流）
   * @param fileId - 文件ID
   */
  serveFile(fileId: string) {
    return apiClient.get<Blob>(`/files/${fileId}`, {
      responseType: "blob",
    });
  },

  /**
   * 获取当前用户的文件列表（分页）
   * @param params - 分页及过滤参数
   * @param params.page - 页码，默认 1
   * @param params.page_size - 每页条数，默认 20
   * @param params.file_type - 按业务类型过滤，可选
   */
  listMyFiles(params?: {
    page?: number;
    page_size?: number;
    file_type?: UploadFileType;
  }) {
    return apiClient.get<ApiResponse<UserFileListResponse>>(
      "/attachments/user/me",
      { params },
    );
  },
};
