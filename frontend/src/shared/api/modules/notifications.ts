/**
 * api/modules/notifications.ts
 */

import { ApiResponse, Notification, PageData } from "@/shared/api/types";
import apiClient from "../client";

export const notificationApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Notification>>>("/notifications", {
      params,
    }),

  unreadCount: () =>
    apiClient.get<ApiResponse<{ count: number }>>(
      "/notifications/unread-count",
    ),

  // 已读所有
  markAllRead: () =>
    apiClient.post<ApiResponse<null>>("/notifications/read-all"),
  // 标记已读
  markRead: (id: number) => apiClient.put(`/notifications/${id}/read`),
};
