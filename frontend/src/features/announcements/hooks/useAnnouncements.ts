// hooks/useAnnouncements.ts
import { useQuery } from '@tanstack/react-query'
import { announcementApi } from '@/shared/api/modules/announcements'
import {
  AnnouncementListParams,
  AnnouncementListResponse,
} from '@/shared/api/types/announcement.model'
import { AnnouncementStatus } from '@/shared/api/types/announcement.model.do'
import { announcementKeys } from './useAnnouncementKeys'

// 辅助：解包后端响应，非 0 code 抛错
function unwrap<T>(res: { code: number; message?: string; data?: T }): T {
  if (res.code !== 0) {
    throw new Error(res.message || '请求失败')
  }
  return res.data as T
}

/**
 * 公告 Hook - 用于前台展示（左侧边栏、公告页面等）
 *
 * 说明：
 * - 所有用户（包括管理员）在前台看到的都是已发布且未过期的公告
 * - 管理员查看所有公告应该在后台管理面板使用 useAdminAnnouncements
 *
 * @param boardId - 可选，板块ID（版主用于获取板块公告）
 */
export function useAnnouncements(boardId?: number) {
  // 使用模块中已定义的类型
  const params: AnnouncementListParams = {
    page: 1,
    page_size: 20,
    status: AnnouncementStatus.Published,
  }

  // 如果传入了板块ID，获取该板块的公告
  if (boardId) {
    params.board_id = boardId
    params.is_global = false
  } else {
    // 默认获取全局公告
    params.is_global = true
  }

  const query = useQuery({
    queryKey: announcementKeys.list(params),
    queryFn: async () => {
      const res = await announcementApi.list(params)
      return unwrap<AnnouncementListResponse>(res.data)
    },
  })

  return {
    announcementsList: query.data?.list || [],
    isLoading: query.isLoading,
    error: query.error
      ? query.error instanceof Error
        ? query.error.message
        : '获取公告失败'
      : null,
    refetch: () => query.refetch().then(() => undefined),
  }
}
