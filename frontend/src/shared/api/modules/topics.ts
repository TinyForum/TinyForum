/**
 * api/modules/topics.ts
 */

import apiClient from "../client";
import type {
  ApiResponse,
  PageData,
  Topic,
  TopicPost,
  TopicFollow,
} from "../types";

export interface CreateTopicPayload {
  title: string;
  description?: string;
  cover?: string;
  is_public?: boolean;
}

export type UpdateTopicPayload = Partial<CreateTopicPayload>;

export const topicApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Topic>>>("/topics", { params }),

  getById: (id: number) => apiClient.get<ApiResponse<Topic>>(`/topics/${id}`),

  getPosts: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<TopicPost>>>(`/topics/${id}/posts`, {
      params,
    }),

  getFollowers: (id: number, params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<TopicFollow>>>(
      `/topics/${id}/followers`,
      { params },
    ),

  create: (data: CreateTopicPayload) =>
    apiClient.post<ApiResponse<Topic>>("/topics", data),

  update: (id: number, data: UpdateTopicPayload) =>
    apiClient.put<ApiResponse<Topic>>(`/topics/${id}`, data),

  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/topics/${id}`),

  addPost: (id: number, data: { post_id: number; sort_order?: number }) =>
    apiClient.post<ApiResponse<null>>(`/topics/${id}/posts`, data),

  removePost: (id: number, postId: number) =>
    apiClient.delete<ApiResponse<null>>(`/topics/${id}/posts/${postId}`),

  follow: (id: number) =>
    apiClient.post<ApiResponse<null>>(`/topics/${id}/follow`),

  unfollow: (id: number) =>
    apiClient.delete<ApiResponse<null>>(`/topics/${id}/follow`),
};
