/**
 * api/modules/notifications.ts
 */

import apiClient from "../client";
import { Notification } from "../types";
import { ApiResponse, PageData } from "../types/basic.model";

export const notificationApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Notification>>>("/notifications", {
      params,
    }),

  unreadCount: () =>
    apiClient.get<ApiResponse<{ count: number }>>(
      "/notifications/count/unread",
    ),

  // 已读所有
  markAllRead: () =>
    apiClient.post<ApiResponse<null>>("/notifications/read-all"),
  // 标记已读
  markRead: (id: number) => apiClient.put(`/notifications/${id}/read`),
};
