// hooks/useNotifications.ts
import { useCallback, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'react-hot-toast'
import { notificationKeys } from './useNotificationKeys'
import { notificationApi } from '@/shared/api/modules/notifications'
import { Notification } from '@/shared/api/types/notification.model'
import { PageData } from '@/shared/api/types/basic.model'

interface UseNotificationsOptions {
  pageSize?: number
  autoLoad?: boolean
}

export function useNotifications(options: UseNotificationsOptions = {}) {
  const { pageSize = 10, autoLoad = true } = options
  const queryClient = useQueryClient()
  const [page, setPage] = useState<number>(1)

  // 查询：分页获取通知列表
  const { data, isLoading, refetch } = useQuery<PageData<Notification>>({
    queryKey: notificationKeys.list({ page, page_size: pageSize }),
    queryFn: async () => {
      const res = await notificationApi.list({ page, page_size: pageSize })
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '获取通知列表失败')
      }
      if (!res.data.data) {
        throw new Error('通知数据为空')
      }
      return res.data.data
    },
    enabled: autoLoad,
  })

  const notifications = data?.list ?? []
  const total = data?.total ?? 0
  // 判断是否还有更多：当前页 * 每页条数 < 总数
  const hasMore = data ? data.page * data.page_size < data.total : true
  const loading = isLoading

  /** 刷新：回到第一页并重新加载 */
  const refresh = useCallback(() => {
    setPage(1)
    refetch()
  }, [refetch])

  /** 加载更多：翻到下一页 */
  const loadMore = useCallback(() => {
    if (!loading && hasMore) {
      setPage((prev) => prev + 1)
    }
  }, [loading, hasMore])

  const markAsReadMutation = useMutation<void, Error, number>({
    mutationFn: async (id: number) => {
      const res = await notificationApi.markRead(id)
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '标记失败')
      }
    },
    onSuccess: () => {
      // 标记已读后刷新通知列表与未读数
      queryClient.invalidateQueries({ queryKey: notificationKeys.lists() })
      queryClient.invalidateQueries({
        queryKey: notificationKeys.unreadCount(),
      })
      queryClient.invalidateQueries({ queryKey: notificationKeys.unread() })
    },
  })

  const markAllAsReadMutation = useMutation<void, Error>({
    mutationFn: async () => {
      const res = await notificationApi.markAllRead()
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '操作失败')
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: notificationKeys.lists() })
      queryClient.invalidateQueries({
        queryKey: notificationKeys.unreadCount(),
      })
      queryClient.invalidateQueries({ queryKey: notificationKeys.unread() })
    },
  })

  /** 标记已读 */
  const markAsRead = async (id: number): Promise<void> => {
    try {
      await markAsReadMutation.mutateAsync(id)
    } catch (err) {
      toast.error(err instanceof Error ? err.message || '操作失败' : '操作失败')
    }
  }

  const markAllAsRead = async (): Promise<void> => {
    try {
      await markAllAsReadMutation.mutateAsync()
    } catch (err) {
      toast.error(err instanceof Error ? err.message || '操作失败' : '操作失败')
    }
  }

  return {
    notifications,
    loading,
    hasMore,
    total,
    refresh,
    loadMore,
    markAsRead,
    markAllAsRead,
  }
}
