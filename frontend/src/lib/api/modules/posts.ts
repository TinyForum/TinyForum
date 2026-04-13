/**
 * api/modules/posts.ts
 * 包含普通帖子 + 问答（question）相关接口
 */

import apiClient from "../client";
import type {
  ApiResponse,
  PageData,
  Post,
  Comment,
  Question,
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
}

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

// ─── 问答 ─────────────────────────────────────────────────────────────────────

export interface QuestionListParams {
  page?: number;
  page_size?: number;
  filter?: "all" | "unanswered" | "answered";
  keyword?: string;
  tag_id?: number;
  board_id?: number;
}

export interface CreateQuestionPayload {
  title: string;
  content: string;
  summary?: string;
  cover?: string;
  board_id?: number;
  tag_ids?: number[];
  reward_score?: number;
}

export interface QuestionDetailResult {
  post: Post;
  liked: boolean;
  question: Question;
  answers: Comment[];
  total: number;
  page: number;
  page_size: number;
}

export interface QuestionDetailParams {
  answer_page?: number;
  answer_page_size?: number;
}

// ─── API ──────────────────────────────────────────────────────────────────────

export const postApi = {
  // ── 普通帖子 ─────────────────────────────────────────────────────────────────
  list: (params?: PostListParams) =>
    apiClient.get<ApiResponse<PageData<Post>>>("/posts", { params }),

  getById: (id: number) =>
    apiClient.get<ApiResponse<PostDetailResult>>(`/posts/${id}`),

  create: (data: CreatePostPayload) =>
    apiClient.post<ApiResponse<Post>>("/posts", data),

  update: (id: number, data: UpdatePostPayload) =>
    apiClient.put<ApiResponse<Post>>(`/posts/${id}`, data),

  delete: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/posts/${id}`),

  like: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/posts/${id}/like`),

  unlike: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/posts/${id}/like`),

  // ── 问答 ─────────────────────────────────────────────────────────────────────
  getQuestions: (params?: QuestionListParams) =>
    apiClient.get<ApiResponse<PageData<Post>>>("/posts/questions", { params }),

  createQuestion: (data: CreateQuestionPayload) =>
    apiClient.post<ApiResponse<Post>>("/posts/question", data),

  getQuestionDetail: (id: number, params?: QuestionDetailParams) =>
    apiClient.get<ApiResponse<QuestionDetailResult>>(
      `/posts/question/${id}`,
      { params }
    ),
  

  /** 采纳答案（通过 post 路由） */
  acceptAnswer: (postId: number, commentId: number) =>
    apiClient.post<ApiResponse<null>>(
      `/posts/questions/${postId}/answer/${commentId}/accept`
    ),

  createAnswer: (postId: number, data: { content: string }) =>
    apiClient.post<ApiResponse<Comment>>(
      `/posts/question/${postId}/answer`,
      data
    ),

  getQuestionAnswers: (
    postId: number,
    params?: { page?: number; page_size?: number }
  ) =>
    apiClient.get<ApiResponse<PageData<Comment>>>(
      `/posts/questions/${postId}/answers`,
      { params }
    ),
};