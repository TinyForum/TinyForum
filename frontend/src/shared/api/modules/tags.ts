/**
 * api/modules/tags.ts
 */

import apiClient from "../client";
import type { ApiResponse, Tag } from "../types";

export interface CreateTagPayload {
  name: string;
  description?: string;
  color?: string;
}

export type UpdateTagPayload = Partial<CreateTagPayload>;

export const tagApi = {
  list: () => apiClient.get<ApiResponse<Tag[]>>("/tags"),

  create: (data: CreateTagPayload) =>
    apiClient.post<ApiResponse<Tag>>("/tags", data),

  update: (id: number, data: UpdateTagPayload) =>
    apiClient.put<ApiResponse<Tag>>(`/tags/${id}`, data),

  delete: (id: number) => apiClient.delete<ApiResponse<null>>(`/tags/${id}`),
};
