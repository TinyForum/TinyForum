// src/lib/api/modules/questions.ts
import apiClient from "../client";
import type {
  ApiResponse,
  PageData,
  QuestionSimple,
  Question,
  Post,
  Comment,
  QuestionResponse,
} from "../types";
import { AnswerListParams } from "./answer";

export interface QuestionListParams {
  page?: number;
  page_size?: number;
  board_id?: number;
  sort?: "latest" | "hot" | "unanswered";
  filter?: string;
  keyword?: string;
}

export const questionApi = {
  // 获取问题精简列表
  getSimple: (params?: QuestionListParams) =>
    apiClient.get<ApiResponse<PageData<QuestionSimple>>>("/questions/simple", {
      params,
    }),

  // 获取问题列表
  getList: (params?: QuestionListParams) =>
    apiClient.get<ApiResponse<PageData<Question>>>("/questions/list", {
      params,
    }),

  // 获取问题详情
  getDetail: (id: number) =>
    apiClient.get<ApiResponse<QuestionResponse>>(`/questions/detail/${id}`),

  // 创建问题
  create: (data: CreateQuestionPayload) =>
    apiClient.post<ApiResponse<Post>>("/questions/create", data),

  // 获取问题的答案列表
  getAnswers: (questionId: number, params?: AnswerListParams) =>
    apiClient.get<ApiResponse<PageData<Comment>>>(
      `/questions/${questionId}/answers`,
      { params },
    ),

  // 创建答案
  createAnswer: (questionId: number, data: { content: string }) =>
    apiClient.post<ApiResponse<Comment>>(
      `/questions/${questionId}/answers`,
      data,
    ),
};

// 创建问题的请求体
export interface CreateQuestionPayload {
  title: string;
  content: string;
  summary?: string;
  cover?: string;
  board_id?: number;
  tag_ids?: number[];
  reward_score?: number;
  post_status?: string;
}
