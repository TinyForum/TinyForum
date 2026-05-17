/**
 * api/modules/comments.ts
 * 包含评论 CRUD + 答案投票 / 采纳
 */

import apiClient from "../client";
import { ApiResponse, PageData } from "../types/basic.model";
import { CreateCommentPayload } from "../types/comment.model";
import { Comment } from "../types/comment.model";

export const commentApi = {
  listByPost: (
    postId: number,
    params?: { page?: number; page_size?: number },
  ) =>
    apiClient.get<ApiResponse<PageData<Comment>>>(`/comments/post/${postId}`, {
      params,
    }),

  create: (data: CreateCommentPayload) =>
    apiClient.post<ApiResponse<Comment>>("/comments", data),

  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/comments/${id}`),

  // ── 答案相关 ──────────────────────────────────────────────────────────────────
  // voteAnswer: (id: number, voteType: VoteType) =>
  //   apiClient.post<ApiResponse<AnswerVoteResult>>(
  //     `/comments/${id}/vote`,
  //     { vote_type: voteType }
  //   ),

  // getVoteStatus: (id: number) =>
  //   apiClient.get<ApiResponse<VoteStatusResult>>(`/comments/${id}/vote`),

  // markAsAnswer: (id: number, isAnswer: boolean) =>
  //   apiClient.put<ApiResponse<null>>(
  //     `/comments/${id}/answer`,
  //     { is_answer: isAnswer }
  // ),
};
