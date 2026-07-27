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
  /**
   * 获取帖子列表
   * @param params
   * @returns
   */
  list: (params?: PostListParams) =>
    apiClient.get<ApiResponse<PageData<Post>>>("/posts", { params }),

  /**
   * 获取帖子详情
   * @param id
   * @returns
   */
  getById: (id: number) =>
    apiClient.get<ApiResponse<PostDetailResult>>(`/posts/${id}`),

  /**
   * 创建文章
   * @param data
   * @returns
   */
  create: (data: CreatePostPayload) =>
    apiClient.post<ApiResponse<Post>>("/posts", data),

  /**
   * 更新文章
   * @param id
   * @param data
   * @returns
   */
  update: (id: number, data: UpdatePostPayload) =>
    apiClient.put<ApiResponse<Post>>(`/posts/${id}`, data),

  /**
   * 删除文章
   * @param id
   * @returns
   */
  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/posts/${id}`),

  /**
   * 点赞
   * @param id
   * @returns
   */
  like: (id: number) => apiClient.post<ApiResponse<null>>(`/posts/${id}/like`),

  /**
   * 取消点赞
   * @param id
   * @returns
   */
  unlike: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/posts/${id}/like`),
};
