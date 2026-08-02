// hooks/useBoard.ts
import { useCallback, useMemo } from 'react'
import { useQuery, UseQueryOptions } from '@tanstack/react-query'
import { boardKeys } from './useBoardKeys'
import { boardApi } from '@/shared/api/modules/boards'
import { Board, BoardPostListItem } from '@/shared/api/types/board.model'
import { PageData } from '@/shared/api/types/basic.model'

// 错误响应类型（用于 retry 判断 404）
interface ErrorResponse {
  response?: {
    status: number
  }
}

interface UseBoardOptions {
  autoLoad?: boolean
  page?: number
  pageSize?: number
}

interface UseBoardReturn {
  boards: Board[]
  loading: boolean
  error: string | null
  total: number
  loadBoards: () => Promise<void>
  getBoardById: (id: number) => Board | undefined
  getDefaultBoard: () => Board | null
  refresh: () => Promise<void>
}

export function useBoard(options: UseBoardOptions = {}): UseBoardReturn {
  const { autoLoad = true, page = 1, pageSize = 100 } = options

  // 查询：分页获取板块列表
  const query = useQuery<PageData<Board>>({
    queryKey: boardKeys.list({ page, page_size: pageSize }),
    queryFn: async () => {
      const res = await boardApi.list({ page, page_size: pageSize })
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '加载板块失败')
      }
      if (!res.data.data) {
        throw new Error('板块数据为空')
      }
      return res.data.data
    },
    enabled: autoLoad,
  })

  const boards = useMemo(() => query.data?.list ?? [], [query.data])
  const total = query.data?.total ?? 0
  const error = query.error ? query.error.message : null

  const loadBoards = useCallback(async () => {
    await query.refetch()
  }, [query])

  const refresh = loadBoards

  const getBoardById = useCallback(
    (id: number): Board | undefined => {
      return boards.find((board: Board): boolean => board.id === id)
    },
    [boards],
  )

  const getDefaultBoard = useCallback((): Board | null => {
    if (boards.length === 0) return null
    return boards[0]
  }, [boards])

  return {
    boards,
    loading: query.isLoading,
    error,
    total,
    loadBoards,
    getBoardById,
    getDefaultBoard,
    refresh,
  }
}

// 通过 slug 获取板块信息
export function useBoardBySlug(
  slug: string,
  options?: Omit<
    UseQueryOptions<Board>,
    'queryKey' | 'queryFn' | 'retry'
  >,
) {
  return useQuery({
    queryKey: boardKeys.slug(slug),
    queryFn: async () => {
      const res = await boardApi.getBySlug(slug)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '加载板块失败')
      }
      if (!res.data.data) {
        throw new Error('板块数据为空')
      }
      return res.data.data
    },
    enabled: !!slug,
    // 404 时不再重试，其余错误最多重试 2 次
    retry: (failureCount: number, error: Error) => {
      const e = error as ErrorResponse
      if (e.response?.status === 404) return false
      return failureCount < 2
    },
    ...options,
  })
}

// 通过 slug 获取板块帖子列表
export function useBoardPostsBySlug(
  slug: string,
  params?: { page?: number; page_size?: number },
  options?: Omit<
    UseQueryOptions<PageData<BoardPostListItem>>,
    'queryKey' | 'queryFn'
  >,
) {
  return useQuery({
    queryKey: boardKeys.postsBySlug(slug, params?.page ?? 1),
    queryFn: async () => {
      const res = await boardApi.getPostsBySlug(slug, params)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '加载板块帖子失败')
      }
      if (!res.data.data) {
        throw new Error('板块帖子数据为空')
      }
      return res.data.data as PageData<BoardPostListItem>
    },
    ...options,
  })
}
