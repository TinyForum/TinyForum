/**
 * api/modules/posts.ts
 * 包含普通帖子 + 问答（question）相关接口
 */

import apiClient from "../client";
import { ApiResponse, PageData } from "../types/basic.model";
import {
  PostListParams,
  Post,
  PostDetailResult,
  CreatePostPayload,
  UpdatePostPayload,
} from "../types/post.model";

// ─── API ──────────────────────────────────────────────────────────────────────

export const postApi = {
  // ── 普通帖子 ─────────────────────────────────────────────────────────────────
  list: (params?: PostListParams) =>
    apiClient.get<ApiResponse<PageData<Post>>>("/posts", { params }),

  getById: (id: number) =>
    apiClient.get<ApiResponse<PostDetailResult>>(`/posts/${id}`),

  /**
   * 创建文章
   * @param data
   * @returns
   */
  create: (data: CreatePostPayload) =>
    apiClient.post<ApiResponse<Post>>("/posts", data),

  update: (id: number, data: UpdatePostPayload) =>
    apiClient.put<ApiResponse<Post>>(`/posts/${id}`, data),

  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/posts/${id}`),

  like: (id: number) => apiClient.post<ApiResponse<null>>(`/posts/${id}/like`),

  unlike: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/posts/${id}/like`),
};
