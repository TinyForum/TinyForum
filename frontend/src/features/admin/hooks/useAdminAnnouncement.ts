import { useCallback } from 'react'
import { useQuery, useQueryClient } from '@tanstack/react-query'
import { announcementApi } from '@/shared/api/modules/announcements'
import toast from 'react-hot-toast'
import { AnnouncementDO } from '@/shared/api/types/announcement.model.do'
import { announcementKeys } from './useAnnouncementKeys'

// 辅助：解包后端响应，非 0 code 抛错
function unwrap<T>(res: { code: number; message?: string; data?: T }): T {
  if (res.code !== 0) {
    throw new Error(res.message || '请求失败')
  }
  return res.data as T
}

// ============ 用于单个公告的 Hook ============
interface UseAdminAnnouncementOptions {
  autoLoad?: boolean
}

interface UseAnnouncementReturn {
  announcement: AnnouncementDO | null
  loading: boolean
  fetch: (id: number) => Promise<AnnouncementDO | null>
  clear: () => void
}

export function useAdminAnnouncement(
  id?: number,
  options?: UseAdminAnnouncementOptions,
): UseAnnouncementReturn {
  const queryClient = useQueryClient()
  const { autoLoad = true } = options || {}

  // 公告详情查询
  const query = useQuery({
    queryKey: announcementKeys.detail(id ?? -1),
    queryFn: async () => {
      const res = await announcementApi.getById(id as number)
      return unwrap<AnnouncementDO>(res.data)
    },
    enabled: autoLoad && !!id,
  })

  // 手动获取公告详情
  const fetch = useCallback(
    async (announcementId: number): Promise<AnnouncementDO | null> => {
      try {
        return await queryClient.fetchQuery({
          queryKey: announcementKeys.detail(announcementId),
          queryFn: async () => {
            const res = await announcementApi.getById(announcementId)
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

  const clear = useCallback(() => {}, [])

  return {
    announcement: query.data ?? null,
    loading: query.isLoading,
    fetch,
    clear,
  }
}
