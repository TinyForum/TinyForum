// hooks/admin/useAdminAnnouncements.ts
import { useCallback, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'
import { announcementApi } from '@/shared/api/modules/announcements'
import { adminAnnouncementApi } from '@/shared/api/modules/admin/announcements'
import {
  AnnouncementListResponse,
  CreateAnnouncementPayload,
  UpdateAnnouncementPayload,
} from '@/shared/api/types/announcement.model'
import {
  AnnouncementDO,
  AnnouncementStatus,
} from '@/shared/api/types/announcement.model.do'
import { announcementKeys } from './useAnnouncementKeys'

// 辅助：解包后端响应，非 0 code 抛错
function unwrap<T>(res: { code: number; message?: string; data?: T }): T {
  if (res.code !== 0) {
    throw new Error(res.message || '请求失败')
  }
  return res.data as T
}

interface UseAdminAnnouncementsReturn {
  // 数据
  announcements: AnnouncementDO[]
  pinnedAnnouncements: AnnouncementDO[]
  total: number

  // 状态
  isLoading: boolean
  isSubmitting: boolean

  // 分页
  page: number
  pageSize: number
  setPage: (page: number) => void
  setPageSize: (pageSize: number) => void

  // 操作方法
  refetch: () => Promise<void>
  getAnnouncementById: (id: number) => Promise<AnnouncementDO | null>
  createAnnouncement: (
    data: CreateAnnouncementPayload,
  ) => Promise<AnnouncementDO | null>
  updateAnnouncement: (
    id: number,
    data: UpdateAnnouncementPayload,
  ) => Promise<AnnouncementDO | null>
  deleteAnnouncement: (id: number) => Promise<boolean>
  publishAnnouncement: (id: number) => Promise<boolean>
  archiveAnnouncement: (id: number) => Promise<boolean>
  pinAnnouncement: (id: number, isPinned: boolean) => Promise<boolean>
}

export function useAdminAnnouncements(): UseAdminAnnouncementsReturn {
  const queryClient = useQueryClient()
  const [page, setPage] = useState<number>(1)
  const [pageSize, setPageSize] = useState<number>(20)

  // 获取公告列表
  const { data, isLoading, refetch: refetchQuery } = useQuery({
    queryKey: announcementKeys.list({ page, page_size: pageSize }),
    queryFn: async () => {
      const res = await announcementApi.adminList({ page, page_size: pageSize })
      return unwrap<AnnouncementListResponse>(res.data)
    },
  })

  const announcements = data?.list || []
  // 筛选置顶公告
  const pinnedAnnouncements = announcements.filter(
    (ann: AnnouncementDO) => ann.is_pinned === true,
  )
  const total = data?.total || 0

  const refetch = useCallback(() => {
    return refetchQuery().then(() => undefined)
  }, [refetchQuery])

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
        toast.error('获取公告详情失败')
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
    mutationFn: async (data) => {
      const res = await adminAnnouncementApi.create(data)
      const result = unwrap<AnnouncementDO>(res.data)
      return result || null
    },
    onSuccess: () => {
      toast.success('创建公告成功')
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
    },
    onError: (error) => {
      console.error('创建公告失败:', error)
      toast.error('创建公告失败')
    },
  })

  const createAnnouncement = useCallback(
    async (data: CreateAnnouncementPayload): Promise<AnnouncementDO | null> => {
      try {
        return await createMutation.mutateAsync(data)
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
    onSuccess: () => {
      toast.success('更新公告成功')
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
    },
    onError: (error) => {
      console.error('更新公告失败:', error)
      toast.error('更新公告失败')
    },
  })

  const updateAnnouncement = useCallback(
    async (
      id: number,
      data: UpdateAnnouncementPayload,
    ): Promise<AnnouncementDO | null> => {
      try {
        return await updateMutation.mutateAsync({ id, data })
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
    },
    onError: (error) => {
      console.error('删除公告失败:', error)
      toast.error('删除公告失败')
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

  // 发布公告：设置 status 为 published，并设置 published_at 为当前时间
  const publishMutation = useMutation<boolean, Error, number>({
    mutationFn: async (id) => {
      const res = await adminAnnouncementApi.update(id, {
        status: AnnouncementStatus.Published,
        published_at: new Date().toISOString(),
      })
      unwrap(res.data)
      return true
    },
    onSuccess: () => {
      toast.success('发布公告成功')
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
    },
    onError: (error) => {
      console.error('发布公告失败:', error)
      toast.error('发布公告失败')
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

  // 归档公告：设置 expired_at 为当前时间之前（标记为过期）
  const archiveMutation = useMutation<boolean, Error, number>({
    mutationFn: async (id) => {
      const res = await adminAnnouncementApi.update(id, {
        expired_at: new Date().toISOString(),
      })
      unwrap(res.data)
      return true
    },
    onSuccess: () => {
      toast.success('归档公告成功')
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
    },
    onError: (error) => {
      console.error('归档公告失败:', error)
      toast.error('归档公告失败')
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
      const res = await adminAnnouncementApi.update(id, { is_pinned: pinned })
      unwrap(res.data)
      return true
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: announcementKeys.lists() })
    },
    onError: (error) => {
      console.error('置顶操作失败:', error)
      toast.error('操作失败')
    },
  })

  const pinAnnouncement = useCallback(
    async (id: number, isPinned: boolean): Promise<boolean> => {
      try {
        return await pinMutation.mutateAsync({ id, pinned: isPinned })
      } catch {
        return false
      }
    },
    [pinMutation],
  )

  const isSubmitting =
    createMutation.isPending ||
    updateMutation.isPending ||
    deleteMutation.isPending ||
    publishMutation.isPending ||
    archiveMutation.isPending ||
    pinMutation.isPending

  return {
    // 数据
    announcements,
    pinnedAnnouncements,
    total,

    // 状态
    isLoading,
    isSubmitting,

    // 分页
    page,
    pageSize,
    setPage,
    setPageSize,

    // 操作方法
    refetch,
    getAnnouncementById,
    createAnnouncement,
    updateAnnouncement,
    deleteAnnouncement,
    publishAnnouncement,
    archiveAnnouncement,
    pinAnnouncement,
  }
}
