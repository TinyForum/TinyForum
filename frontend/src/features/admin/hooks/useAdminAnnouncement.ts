import { useCallback, useEffect, useState } from 'react'
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
  const [activeId, setActiveId] = useState<number | null>(id ?? null)

  // 同步外部传入的 id
  useEffect(() => {
    if (id !== undefined) {
      setActiveId(id)
    }
  }, [id])

  // 公告详情查询
  const query = useQuery({
    queryKey: announcementKeys.detail(activeId ?? -1),
    queryFn: async () => {
      const res = await announcementApi.getById(activeId as number)
      return unwrap<AnnouncementDO>(res.data)
    },
    enabled: autoLoad && !!activeId,
  })

  // 手动获取公告详情
  const fetch = useCallback(
    async (announcementId: number): Promise<AnnouncementDO | null> => {
      setActiveId(announcementId)
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

  const clear = useCallback(() => {
    setActiveId(null)
  }, [])

  return {
    announcement: query.data ?? null,
    loading: query.isLoading,
    fetch,
    clear,
  }
}
