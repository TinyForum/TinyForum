// src/lib/api/modules/questions.ts
import apiClient from "../client";
import type { ApiResponse, PageData, QuestionSimple, Question, Post } from "../types";

export interface QuestionListParams {
  page?: number;
  page_size?: number;
  board_id?: number;
  sort?: 'latest' | 'hot' | 'unanswered';
}

export const questionApi = {
  // 获取问题精简列表
  getSimple: (params?: QuestionListParams) =>
    apiClient.get<ApiResponse<PageData<QuestionSimple>>>("/questions/simple", { params }),
  
  // 获取问题详情
  getDetail: (id: number) =>
    apiClient.get<ApiResponse<Question>>(`/questions/${id}`),
  
  // 创建问题
  create: (data: CreateQuestionPayload) =>
    apiClient.post<ApiResponse<Post>>("/questions", data),
  
  // 创建回答
  createAnswer: (questionId: number, data: { content: string }) =>
    apiClient.post<ApiResponse<Comment>>(`/questions/${questionId}/answers`, data),
  
  // 采纳答案
  acceptAnswer: (questionId: number, answerId: number) =>
    apiClient.post<ApiResponse<null>>(`/questions/${questionId}/answers/${answerId}/accept`),
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
}