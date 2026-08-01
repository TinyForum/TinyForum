/**
 * api/modules/comments.ts
 * 包含评论 CRUD + 答案投票 / 采纳
 */

import apiClient from "../client";
import { ApiResponse, PageData } from "../types/basic.model";
import {
  CreateCommentPayload,
  CreateCommentResponse,
} from "../types/comment.model";
import { CommentResponse } from "../types/comment.model";

export const commentApi = {
  /**
   * list comments
   * @param postId
   * @param params
   * @returns
   */
  listByPost: (
    postId: number,
    params?: { page?: number; page_size?: number },
  ) =>
    apiClient.get<ApiResponse<PageData<CommentResponse>>>(
      `/comments/post/${postId}`,
      {
        params,
      },
    ),

  /**
   * 获取评论树
   * @param postId
   * @returns
   */
  listTree: (postId: number) =>
    apiClient.get<ApiResponse<CommentResponse[]>>(
      `/comments/post/${postId}/tree`,
    ),
  /**
   * 创建评论
   */
  create: (data: CreateCommentPayload) =>
    apiClient.post<ApiResponse<CreateCommentResponse>>("/comments", data),

  /**
   * delete the comment
   * @param id
   * @returns
   */
  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/comments/${id}`),
};
