// hooks/useBoard.ts
import { useCallback, useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { boardKeys } from './useBoardKeys'
import { boardApi } from '@/shared/api/modules/boards'
import { Board } from '@/shared/api/types/board.model'
import { PageData } from '@/shared/api/types/basic.model'

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
