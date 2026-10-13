import apiClient from "../client";
import { ApiResponse, PageData } from "../types/basic.model";
import { AuditTaskVO } from "../types/reviewer.model";

export interface ReviewTaskPayload {
  note?: string;
}

export const reviewerApi = {
  // 待审核任务列表
  listPendingTasks: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<AuditTaskVO>>>("/reviewer/audit/tasks", {
      params,
    }),
  // 审核历史（已处理任务）
  listHistory: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<AuditTaskVO>>>("/reviewer/audit/history", {
      params,
    }),
  // 审核通过
  approveTask: (taskId: number, data?: ReviewTaskPayload) =>
    apiClient.put<ApiResponse<{ reviewed: boolean }>>(
      `/reviewer/audit/tasks/${taskId}/approve`,
      data ?? {},
    ),
  // 审核拒绝
  rejectTask: (taskId: number, data?: ReviewTaskPayload) =>
    apiClient.put<ApiResponse<{ reviewed: boolean }>>(
      `/reviewer/audit/tasks/${taskId}/reject`,
      data ?? {},
    ),
};
