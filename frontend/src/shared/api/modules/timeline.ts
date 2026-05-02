/**
 * api/modules/timeline.ts
 */

import apiClient from "../client";
import type { TimelineEvent, Subscription } from "../types";
import { ApiResponse, PageData } from "../types/basic.model";

export const timelineApi = {
  getHome: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<TimelineEvent>>>("/timeline", {
      params,
    }),

  getFollowing: (params?: { page?: number; page_size?: number }) =>
    apiClient.get<ApiResponse<PageData<TimelineEvent>>>("/timeline/following", {
      params,
    }),

  subscribe: (userId: number) =>
    apiClient.post<ApiResponse<null>>(`/timeline/subscribe/${userId}`),

  unsubscribe: (userId: number) =>
    apiClient.delete<ApiResponse<null>>(`/timeline/subscribe/${userId}`),

  getSubscriptions: () =>
    apiClient.get<ApiResponse<Subscription[]>>("/timeline/subscriptions"),
};
