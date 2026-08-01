// hooks/useFollowList.ts
import { useCallback, useState } from "react"
import { useQuery } from "@tanstack/react-query"
import { userApi } from "@/shared/api/modules/user"
import { UserDO } from "@/shared/api/types/user.model.do"
import { PageData } from "@/shared/api/types/basic.model"
import { userKeys } from "./useUserKeys"

// 当前列表请求参数
interface FollowListQuery {
  userId: number | null
  kind: "followers" | "following"
  page: number
  pageSize: number
}

// ========== 粉丝/关注列表（分页） ==========
interface UseFollowListReturn {
  users: UserDO[]
  loading: boolean
  error: string | null
  total: number
  page: number
  loadFollowers: (
    userId: number,
    pageNum?: number,
    pageSize?: number,
  ) => Promise<void>
  loadFollowing: (
    userId: number,
    pageNum?: number,
    pageSize?: number,
  ) => Promise<void>
}

export function useFollowList(): UseFollowListReturn {
  const [query, setQuery] = useState<FollowListQuery>({
    userId: null,
    kind: "followers",
    page: 1,
    pageSize: 20,
  })

  // 查询：按 userId + kind 拉取粉丝或关注列表
  const { data, isLoading, error } = useQuery<PageData<UserDO>>({
    queryKey:
      query.kind === "followers"
        ? userKeys.followers(query.userId ?? 0, {
            page: query.page,
            page_size: query.pageSize,
          })
        : userKeys.following(query.userId ?? 0, {
            page: query.page,
            page_size: query.pageSize,
          }),
    queryFn: async () => {
      if (query.userId === null) {
        throw new Error("缺少用户 ID")
      }
      const params = { page: query.page, page_size: query.pageSize }
      const res =
        query.kind === "followers"
          ? await userApi.getFollowers(query.userId, params)
          : await userApi.getFollowing(query.userId, params)
      const failMsg =
        query.kind === "followers" ? "获取粉丝列表失败" : "获取关注列表失败"
      if (res.data.code !== 0) {
        throw new Error(res.data.message || failMsg)
      }
      if (!res.data.data) {
        throw new Error(failMsg)
      }
      return res.data.data
    },
    enabled: query.userId !== null,
  })

  // 命令式入口：仅更新本地查询参数，由 useQuery 自动触发请求
  const loadFollowers = useCallback(
    async (
      userId: number,
      pageNum: number = 1,
      pageSize: number = 20,
    ): Promise<void> => {
      setQuery({ userId, kind: "followers", page: pageNum, pageSize })
    },
    [],
  )

  const loadFollowing = useCallback(
    async (
      userId: number,
      pageNum: number = 1,
      pageSize: number = 20,
    ): Promise<void> => {
      setQuery({ userId, kind: "following", page: pageNum, pageSize })
    },
    [],
  )

  return {
    users: data?.list ?? [],
    loading: isLoading,
    error: error?.message ?? null,
    total: data?.total ?? 0,
    page: data?.page ?? 1,
    loadFollowers,
    loadFollowing,
  }
}
