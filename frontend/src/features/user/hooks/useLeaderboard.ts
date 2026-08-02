import { useCallback, useState } from "react"
import { useQuery } from "@tanstack/react-query"
import { userApi } from "@/shared/api/modules/user"
import { LeaderboardItemResponse } from "@/shared/api/types/user.model"
import { userKeys } from "./useUserKeys"

// ========== 排行榜 ==========
interface UseLeaderboardReturn {
  data: LeaderboardItemResponse[]
  loading: boolean
  error: string | null
  loadLeaderboard: (simple?: boolean, limit?: number) => Promise<void>
}

interface LeaderboardParams {
  simple: boolean
  limit: number
}

export function useLeaderboard(): UseLeaderboardReturn {
  const [params, setParams] = useState<LeaderboardParams>({
    simple: true,
    limit: 100,
  })

  // 查询：按参数拉取排行榜
  const { data, isLoading, error } = useQuery<LeaderboardItemResponse[]>({
    queryKey: userKeys.leaderboard(params.simple, params.limit),
    queryFn: async () => {
      const res = params.simple
        ? await userApi.getLeaderboardSimple({ limit: params.limit })
        : await userApi.getLeaderboardDetail({ limit: params.limit })
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "加载排行榜失败")
      }
      if (!res.data.data) {
        throw new Error("加载排行榜失败")
      }
      return res.data.data
    },
  })

  // 命令式入口：更新查询参数
  const loadLeaderboard = useCallback(
    async (simple: boolean = true, limit: number = 100): Promise<void> => {
      setParams({ simple, limit })
    },
    [],
  )

  return {
    data: data ?? [],
    loading: isLoading,
    error: error?.message ?? null,
    loadLeaderboard,
  }
}
