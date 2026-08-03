import apiClient from "../client";
import { ApiResponse } from "../types/basic.model";
import {
  RecommendationResponse,
  BehaviorEvent,
  RecommendationFeedback,
  BatchFeedback,
  UserInterestProfile,
} from "../types/recommendation.model";

export const recommendationApi = {
  getFeed: (params?: { page?: number; page_size?: number; strategy?: string }) =>
    apiClient.get<ApiResponse<RecommendationResponse>>("/recommendations/feed", {
      params,
    }),

  recordBehavior: (data: BehaviorEvent) =>
    apiClient.post<ApiResponse<null>>("/recommendations/behaviors", data),

  recordBehaviorsBatch: (data: BehaviorEvent[]) =>
    apiClient.post<ApiResponse<null>>("/recommendations/behaviors/batch", data),

  submitFeedback: (data: RecommendationFeedback) =>
    apiClient.post<ApiResponse<null>>("/recommendations/feedback", data),

  submitFeedbackBatch: (data: BatchFeedback) =>
    apiClient.post<ApiResponse<null>>("/recommendations/feedback/batch", data),

  getProfile: () =>
    apiClient.get<ApiResponse<UserInterestProfile>>("/recommendations/profile"),

  refreshProfile: () =>
    apiClient.post<ApiResponse<null>>("/recommendations/profile/refresh"),

  recordView: (data: BehaviorEvent) =>
    apiClient.post<ApiResponse<null>>("/recommendations/record-view", data),
};
