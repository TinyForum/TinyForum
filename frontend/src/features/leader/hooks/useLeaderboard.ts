// hooks/useLeaderboard.ts
import {
  LeaderboardRequest,
  LeaderboardItemResponse,
  userAPI,
} from "@/shared/api/modules/users";
import { useQuery, UseQueryOptions } from "@tanstack/react-query";

/**
 * 获取排行榜数据
 * @param params 查询参数（limit, fields）
 * @param options React Query 配置选项
 */
export const useLeaderboard = (
  params?: LeaderboardRequest,
  options?: Omit<
    UseQueryOptions<LeaderboardItemResponse[], Error>,
    "queryKey" | "queryFn"
  >,
) => {
  return useQuery({
    queryKey: ["leaderboard", params?.limit],
    queryFn: async (): Promise<LeaderboardItemResponse[]> => {
      const { data } = await userAPI.getLeaderboardDetail(params);
      // 确保返回数组，如果 data.data 为 undefined 则返回空数组
      return data.data || [];
    },
    staleTime: 5 * 60 * 1000, // 5 分钟内数据视为新鲜
    ...options,
  });
};
