import { useCallback } from "react"
import { useQuery } from "@tanstack/react-query"
import { userStatsApi } from "@/shared/api/modules/user/stats"
import { UserStatsVO } from "@/shared/api/types/user.model"
import { userKeys } from "./useUserKeys"

type UseUserStatsReturn = UserStatsVO & {
  loadStats: () => Promise<void>
  isLoading: boolean
  error: string | null
}

export function useUserStats(): UseUserStatsReturn {
  // 查询：拉取当前用户统计数据
  const { data, isLoading, error, refetch } = useQuery<UserStatsVO>({
    queryKey: userKeys.stats(),
    queryFn: async () => {
      const res = await userStatsApi.getUserStats()
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取统计数量失败")
      }
      if (!res.data.data) {
        throw new Error("获取统计数量失败")
      }
      return res.data.data
    },
  })

  // 命令式入口：重新拉取统计
  const loadStats = useCallback(async (): Promise<void> => {
    await refetch()
  }, [refetch])

  return {
    // 返回值与 UserStatsVO 完全对齐（下划线命名）
    total_post: data?.total_post ?? 0,
    total_comment: data?.total_comment ?? 0,
    total_favorite: data?.total_favorite ?? 0,
    total_like: data?.total_like ?? 0,
    total_follower: data?.total_follower ?? 0,
    total_following: data?.total_following ?? 0,
    total_report: data?.total_report ?? 0,
    total_violation: data?.total_violation ?? 0,
    total_question: data?.total_question ?? 0,
    total_answer: data?.total_answer ?? 0,
    total_upload: data?.total_upload ?? 0,
    total_score: data?.total_score ?? 0,
    isLoading,
    error: error?.message ?? null,
    loadStats,
  }
}
