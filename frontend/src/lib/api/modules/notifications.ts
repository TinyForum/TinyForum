/**
 * api/modules/notifications.ts
 */

import apiClient from "../client";
import type { ApiResponse, PageData, Notification } from "../types";

export const notificationApi = {
  list: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<Notification>>>("/notifications", {
      params,
    }),

  unreadCount: () =>
    apiClient.get<ApiResponse<{ count: number }>>(
      "/notifications/unread-count"
    ),

  markAllRead: () =>
    apiClient.post<ApiResponse<null>>("/notifications/read-all"),
};