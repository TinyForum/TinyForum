/**
 * api/modules/posts.ts
 * 包含普通帖子 + 问答（question）相关接口
 */

import apiClient from "../client";
import type {
  ApiResponse,
  PageData,
  Post,
  // Comment,
  // Question,
  PostType,
} from "../types";

// ─── 普通帖子 ─────────────────────────────────────────────────────────────────

export interface PostListParams {
  page?: number;
  page_size?: number;
  keyword?: string;
  sort_by?: string;
  type?: PostType;
  author_id?: number;
  tag_id?: number;
  board_id?: number;
}

export interface CreatePostPayload {
  title: string;
  content: string;
  summary?: string;
  cover?: string;
  type?: PostType;
  board_id?: number;
  tag_ids?: number[];
  status?: PostStatus;
}

export type PostStatus = "draft" | "published" | "pending" | "hidden";

export interface UpdatePostPayload {
  title?: string;
  content?: string;
  summary?: string;
  cover?: string;
  tag_ids?: number[];
}

export interface PostDetailResult {
  post: Post;
  liked: boolean;
}

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
