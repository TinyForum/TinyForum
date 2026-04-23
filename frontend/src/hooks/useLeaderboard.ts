// hooks/useLeaderboard.ts
import { LeaderboardRequest, LeaderboardItemResponse, userApi } from '@/lib/api/modules/users';
import { useQuery, UseQueryOptions } from '@tanstack/react-query';
// import { userApi } from '@/services/user';
// import { LeaderboardRequest, LeaderboardResponse } from '@/types/user';

/**
 * 获取排行榜数据
 * @param params 查询参数（limit, fields）
 * @param options React Query 配置选项
 */
export const useLeaderboard = (
  params?: LeaderboardRequest,
  options?: Omit<UseQueryOptions<LeaderboardItemResponse[], Error>, 'queryKey' | 'queryFn'>
) => {
  return useQuery({
    queryKey: ['leaderboard', params?.limit, params?.fields],
    queryFn: async () => {
      const { data } = await userApi.getLeaderboardSimple(params);
      return data.data; 
    },
    staleTime: 5 * 60 * 1000, // 5 分钟内数据视为新鲜
    ...options,
  });
};