import { useState, useCallback } from "react"
import { useQuery } from "@tanstack/react-query"
import { userPostApi } from "@/shared/api/modules/user/post"
import {
  UserPostsVO,
  GetUserPostsRequest,
} from "@/shared/api/types/user.model"
import { PageData } from "@/shared/api/types/basic.model"
import { userKeys } from "./useUserKeys"

// 分页数据核心字段（复用后端 PageData 结构）
type UserPostsPageData = PageData<UserPostsVO>

// Hook 返回类型：分页数据 + 控制字段
type UseMePostsReturn = UserPostsPageData & {
  isLoading: boolean
  error: string | null
  loadPosts: (params?: Partial<GetUserPostsRequest>) => Promise<void>
  loadMore: () => Promise<void>
  refresh: () => Promise<void>
}

// 默认请求参数（避免 undefined 传递给 API）
const DEFAULT_PARAMS: GetUserPostsRequest = {
  page: 1,
  page_size: 10,
}

export function useMePosts(
  initialParams?: Partial<GetUserPostsRequest>,
): UseMePostsReturn {
  // 当前完整查询参数（合并默认值）
  const [params, setParams] = useState<GetUserPostsRequest>({
    ...DEFAULT_PARAMS,
    ...initialParams,
  })

  // 查询：按当前参数拉取当前用户帖子列表
  const { data, isFetching, error, refetch } = useQuery<PageData<UserPostsVO>>({
    queryKey: userKeys.list(params),
    queryFn: async () => {
      const res = await userPostApi.listUserPosts(params)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || "获取帖子列表失败")
      }
      if (!res.data.data) {
        throw new Error("返回数据为空")
      }
      return res.data.data
    },
  })

  const list = data?.list ?? []
  const total = data?.total ?? 0
  const page = data?.page ?? params.page
  const pageSize = data?.page_size ?? params.page_size
  const hasMore = data?.has_more ?? false

  // 命令式入口：合并参数后触发查询
  const loadPosts = useCallback(
    async (next?: Partial<GetUserPostsRequest>): Promise<void> => {
      setParams((prev) => ({
        ...prev,
        ...next,
        page: next?.page ?? prev.page,
        page_size: next?.page_size ?? prev.page_size,
      }))
    },
    [],
  )

  // 加载下一页（基于当前 page+1）
  const loadMore = useCallback(async (): Promise<void> => {
    if (isFetching || !hasMore) return
    setParams((prev) => ({ ...prev, page: prev.page + 1 }))
  }, [isFetching, hasMore])

  // 刷新（重新加载第一页）
  const refresh = useCallback(async (): Promise<void> => {
    setParams((prev) => ({ ...prev, page: 1 }))
    await refetch()
  }, [refetch])

  return {
    list,
    total,
    page,
    page_size: pageSize,
    has_more: hasMore,
    isLoading: isFetching,
    error: error?.message ?? null,
    loadPosts,
    loadMore,
    refresh,
  }
}
