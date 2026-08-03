import { useQuery, useMutation, useQueryClient, UseQueryOptions } from "@tanstack/react-query";
import { recommendationApi } from "@/shared/api/modules/recommendation";
import { RecommendationResponse } from "@/shared/api/types/recommendation.model";

export const recommendationKeys = {
  all: ["recommendations"] as const,
  feed: (params?: { page?: number; page_size?: number; strategy?: string }) =>
    [...recommendationKeys.all, "feed", params] as const,
  profile: () => [...recommendationKeys.all, "profile"] as const,
};

export function useRecommendationFeed(
  params?: { page?: number; page_size?: number; strategy?: string },
  options?: Partial<UseQueryOptions<RecommendationResponse>>,
) {
  return useQuery<RecommendationResponse>({
    queryKey: recommendationKeys.feed(params),
    queryFn: async () => {
      const response = await recommendationApi.getFeed(params);
      if (response.data.code !== 0) {
        throw new Error(response.data.message || "Failed to fetch recommendations");
      }
      return response.data.data as RecommendationResponse;
    },
    staleTime: 30_000,
    ...options,
  });
}

export function useRecordBehavior() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: { target_id: number; target_type: string; behavior_type: string }) =>
      recommendationApi.recordBehavior(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: recommendationKeys.all });
    },
  });
}

export function useRecordView() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: { target_id: number; target_type: string }) =>
      recommendationApi.recordView({ ...data, behavior_type: "view" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: recommendationKeys.all });
    },
  });
}

export function useUserInterestProfile() {
  return useQuery({
    queryKey: recommendationKeys.profile(),
    queryFn: async () => {
      const response = await recommendationApi.getProfile();
      if (response.data.code !== 0) {
        throw new Error(response.data.message || "Failed to fetch profile");
      }
      return response.data.data;
    },
  });
}

export function useRefreshProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => recommendationApi.refreshProfile(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: recommendationKeys.profile() });
    },
  });
}
