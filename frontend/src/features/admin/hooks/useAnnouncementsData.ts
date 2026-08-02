// hooks/admin/useAnnouncementsData.ts
import { useCallback, useMemo, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { adminAnnouncementApi } from '@/shared/api/modules/admin/announcements'
import { announcementApi } from '@/shared/api/modules/announcements'
import {
  AnnouncementListParams,
  AnnouncementListResponse,
  CreateAnnouncementPayload,
  UpdateAnnouncementPayload,
} from '@/shared/api/types/announcement.model'
import { AnnouncementDO } from '@/shared/api/types/announcement.model.do'
import { announcementKeys } from './useAnnouncementKeys'

// 辅助：解包后端响应，非 0 code 抛错
function unwrap<T>(res: { code: number; message?: string; data?: T }): T {
  if (res.code !== 0) {
    throw new Error(res.message || '请求失败')
  }
  return res.data as T
}

// ============ 配置选项 ============
interface UseAnnouncementsDataOptions {
  enabled?: boolean
  defaultParams?: AnnouncementListParams
  autoLoadPinned?: boolean
}

// ============ Hook 返回值类型 ============
interface UseAnnouncementsDataReturn {
  // 数据状态
  announcements: AnnouncementDO[]
  pinnedAnnouncements: AnnouncementDO[]
  total: number
  loading: boolean
  submitting: boolean
  refreshing: boolean
  isLoading: boolean

  // 分页
  page: number
  pageSize: number

  // 操作方法
  fetchAnnouncements: (params?: AnnouncementListParams) => Promise<void>
  fetchPinnedAnnouncements: (boardId?: number) => Promise<void>
  getAnnouncementById: (id: number) => Promise<AnnouncementDO | null>
  createAnnouncement: (
    params: CreateAnnouncementPayload,
  ) => Promise<AnnouncementDO | null>
  updateAnnouncement: (
    id: number,
    params: UpdateAnnouncementPayload,
  ) => Promise<AnnouncementDO | null>
  deleteAnnouncement: (id: number) => Promise<boolean>
  publishAnnouncement: (id: number) => Promise<boolean>
  archiveAnnouncement: (id: number) => Promise<boolean>
  pinAnnouncement: (id: number, pinned: boolean) => Promise<boolean>

  // 状态设置
  setPage: (page: number) => void
  setPageSize: (pageSize: number) => void
  setFilters: (filters: AnnouncementListParams) => void
  resetFilters: () => void
}

// ============ Hook 实现 ============
export function useAnnouncementsData(
  enabled: boolean = true,
  options?: UseAnnouncementsDataOptions,
): UseAnnouncementsDataReturn {
  const queryClient = useQueryClient()
  const { defaultParams = {}, autoLoadPinned = true } = options || {}

  // 分页与筛选状态
  const [page, setPage] = useState<number>(defaultParams.page || 1)
  const [pageSize, setPageSize] = useState<number>(
    defaultParams.page_size || 20,
  )
  const [filters, setFilters] = useState<AnnouncementListParams>(() => ({
    page,
    page_size: pageSize,
    ...defaultParams,
  }))
  const [pinnedBoardId, setPinnedBoardId] = useState<number | undefined>(
    undefined,
  )

  // 合并筛选条件与分页参数，并移除 undefined 值
  const listParams = useMemo(() => {
    const params: AnnouncementListParams = {
      ...filters,
      page,
      page_size: pageSize,
    }
    Object.keys(params).forEach((key) => {
      if (params[key as keyof AnnouncementListParams] === undefined) {
        delete params[key as keyof AnnouncementListParams]
      }
    })
    return params
  }, [filters, page, pageSize])

  // 获取公告列表
  const listQuery = useQuery({
    queryKey: announcementKeys.list(listParams),
    queryFn: async () => {
      const res = await announcementApi.adminList(listParams)
      return unwrap<AnnouncementListResponse>(res.data)
    },
    enabled,
  })

  // 获取置顶公告
  const pinnedQuery = useQuery({
    queryKey: announcementKeys.pinned(pinnedBoardId),
    queryFn: async () => {
      const res = await announcementApi.getPinned(pinnedBoardId)
      return unwrap<AnnouncementDO[]>(res.data)
    },
    enabled: enabled && autoLoadPinned,
  })

  const announcements = listQuery.data?.list || []
  const pinnedAnnouncements = pinnedQuery.data || []
  const total = listQuery.data?.total || 0

  // 刷新公告列表
  const fetchAnnouncements = useCallback(async (): Promise<void> => {
    if (!enabled) return
    await listQuery.refetch()
  }, [enabled, listQuery])

  // 刷新置顶公告
  const fetchPinnedAnnouncements = useCallback(
    async (boardId?: number): Promise<void> => {
      if (!enabled) return
      setPinnedBoardId(boardId)
      await queryClient.fetchQuery({
        queryKey: announcementKeys.pinned(boardId),
        queryFn: async () => {
          const res = await announcementApi.getPinned(boardId)
          return unwrap<AnnouncementDO[]>(res.data)
        },
      })
    },
    [enabled, queryClient],
  )

  // 根据 ID 获取公告详情
  const getAnnouncementById = useCallback(
    async (id: number): Promise<AnnouncementDO | null> => {
      try {
        return await queryClient.fetchQuery({
          queryKey: announcementKeys.detail(id),
          queryFn: async () => {
            const res = await announcementApi.getById(id)
            return unwrap<AnnouncementDO>(res.data)
          },
        })
      } catch (error) {
        console.error('获取公告详情失败:', error)
        toast.error('获取公告详情失败，请稍后重试')
        return null
      }
    },
    [queryClient],
  )

  // 创建公告
  const createMutation = useMutation<
    AnnouncementDO | null,
    Error,
    CreateAnnouncementPayload
  >({
    mutationFn: async (params) => {
      const res = await adminAnnouncementApi.create(params)
      const result = unwrap<AnnouncementDO>(res.data)
      return result || null
    },
    onSuccess: () => {
      toast.success('创建公告成功')
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
    },
    onError: (error) => {
      console.error('创建公告失败:', error)
      toast.error('创建公告失败，请稍后重试')
    },
  })

  const createAnnouncement = useCallback(
    async (params: CreateAnnouncementPayload): Promise<AnnouncementDO | null> => {
      try {
        return await createMutation.mutateAsync(params)
      } catch {
        return null
      }
    },
    [createMutation],
  )

  // 更新公告
  const updateMutation = useMutation<
    AnnouncementDO | null,
    Error,
    { id: number; data: UpdateAnnouncementPayload }
  >({
    mutationFn: async ({ id, data }) => {
      const res = await adminAnnouncementApi.update(id, data)
      const result = unwrap<AnnouncementDO>(res.data)
      return result || null
    },
    onSuccess: (_, { data }) => {
      toast.success('更新公告成功')
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
      // 如果影响置顶状态，刷新置顶列表
      if (data.is_pinned !== undefined) {
        queryClient.invalidateQueries({ queryKey: announcementKeys.pinned() })
      }
    },
    onError: (error) => {
      console.error('更新公告失败:', error)
      toast.error('更新公告失败，请稍后重试')
    },
  })

  const updateAnnouncement = useCallback(
    async (
      id: number,
      params: UpdateAnnouncementPayload,
    ): Promise<AnnouncementDO | null> => {
      try {
        return await updateMutation.mutateAsync({ id, data: params })
      } catch {
        return null
      }
    },
    [updateMutation],
  )

  // 删除公告
  const deleteMutation = useMutation<boolean, Error, number>({
    mutationFn: async (id) => {
      const res = await adminAnnouncementApi.delete(id)
      unwrap(res.data)
      return true
    },
    onSuccess: () => {
      toast.success('删除公告成功')
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
      queryClient.invalidateQueries({ queryKey: announcementKeys.pinned() })
    },
    onError: (error) => {
      console.error('删除公告失败:', error)
      toast.error('删除公告失败，请稍后重试')
    },
  })

  const deleteAnnouncement = useCallback(
    async (id: number): Promise<boolean> => {
      try {
        return await deleteMutation.mutateAsync(id)
      } catch {
        return false
      }
    },
    [deleteMutation],
  )

  // 发布公告
  const publishMutation = useMutation<boolean, Error, number>({
    mutationFn: async (id) => {
      const res = await adminAnnouncementApi.publish(id)
      unwrap(res.data)
      return true
    },
    onSuccess: () => {
      toast.success('发布公告成功')
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
      queryClient.invalidateQueries({ queryKey: announcementKeys.pinned() })
    },
    onError: (error) => {
      console.error('发布公告失败:', error)
      toast.error('发布公告失败，请稍后重试')
    },
  })

  const publishAnnouncement = useCallback(
    async (id: number): Promise<boolean> => {
      try {
        return await publishMutation.mutateAsync(id)
      } catch {
        return false
      }
    },
    [publishMutation],
  )

  // 归档公告
  const archiveMutation = useMutation<boolean, Error, number>({
    mutationFn: async (id) => {
      const res = await adminAnnouncementApi.archive(id)
      unwrap(res.data)
      return true
    },
    onSuccess: () => {
      toast.success('归档公告成功')
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
    },
    onError: (error) => {
      console.error('归档公告失败:', error)
      toast.error('归档公告失败，请稍后重试')
    },
  })

  const archiveAnnouncement = useCallback(
    async (id: number): Promise<boolean> => {
      try {
        return await archiveMutation.mutateAsync(id)
      } catch {
        return false
      }
    },
    [archiveMutation],
  )

  // 置顶/取消置顶
  const pinMutation = useMutation<
    boolean,
    Error,
    { id: number; pinned: boolean }
  >({
    mutationFn: async ({ id, pinned }) => {
      const res = await adminAnnouncementApi.pin(id, pinned)
      unwrap(res.data)
      return true
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
      queryClient.invalidateQueries({ queryKey: announcementKeys.pinned() })
    },
    onError: (error) => {
      console.error('置顶操作失败:', error)
      toast.error('操作失败，请稍后重试')
    },
  })

  const pinAnnouncement = useCallback(
    async (id: number, pinned: boolean): Promise<boolean> => {
      try {
        return await pinMutation.mutateAsync({ id, pinned })
      } catch {
        return false
      }
    },
    [pinMutation],
  )

  // 重置筛选条件
  const resetFilters = useCallback(() => {
    const newFilters = {
      page: 1,
      page_size: pageSize,
    }
    setFilters(newFilters)
    setPage(1)
  }, [pageSize])

  const submitting =
    createMutation.isPending ||
    updateMutation.isPending ||
    deleteMutation.isPending ||
    publishMutation.isPending ||
    archiveMutation.isPending ||
    pinMutation.isPending

  return {
    // 数据状态
    announcements,
    pinnedAnnouncements,
    total,
    loading: listQuery.isFetching,
    submitting,
    refreshing: listQuery.isRefetching,
    isLoading: listQuery.isLoading,

    // 分页
    page,
    pageSize,

    // 操作方法
    fetchAnnouncements,
    fetchPinnedAnnouncements,
    getAnnouncementById,
    createAnnouncement,
    updateAnnouncement,
    deleteAnnouncement,
    publishAnnouncement,
    archiveAnnouncement,
    pinAnnouncement,

    // 状态设置
    setPage,
    setPageSize,
    setFilters,
    resetFilters,
  }
}
