// hooks/useUnreadCount.ts
import { useState, useEffect, useCallback } from 'react'
import { useQuery } from '@tanstack/react-query'
import { notificationKeys } from './useNotificationKeys'
import { notificationApi } from '@/shared/api/modules/notifications'

export function useUnreadCount(isAuthenticated: boolean = false) {
  const [unreadCount, setUnreadCount] = useState<number>(0)

  // 查询：获取未读数量
  const query = useQuery<{ count: number }>({
    queryKey: notificationKeys.unreadCount(),
    queryFn: async () => {
      const res = await notificationApi.unreadCount()
      if (res.data.code !== 0) {
        throw new Error(res.data.message || '获取未读数量失败')
      }
      if (!res.data.data) {
        throw new Error('未读数量数据为空')
      }
      return res.data.data
    },
    enabled: isAuthenticated,
  })

  // 查询结果同步到本地状态，保证与服务端数据一致
  useEffect(() => {
    if (query.data) {
      setUnreadCount(query.data.count)
    }
  }, [query.data])

  // 手动增加/减少计数（用于本地乐观更新）
  const incrementUnread = useCallback(() => {
    setUnreadCount((prev) => prev + 1)
  }, [])

  const decrementUnread = useCallback(() => {
    setUnreadCount((prev) => Math.max(0, prev - 1))
  }, [])

  const resetUnread = useCallback(() => {
    setUnreadCount(0)
  }, [])

  const refetch = useCallback(async () => {
    const result = await query.refetch()
    if (result.data) {
      setUnreadCount(result.data.count)
    }
  }, [query])

  return {
    unreadCount,
    loading: query.isLoading,
    refetch,
    incrementUnread,
    decrementUnread,
    resetUnread,
  }
}
