// hooks/useNotificationKeys.ts
export const notificationKeys = {
  all: ['notifications'] as const,
  lists: () => [...notificationKeys.all, 'list'] as const,
  list: (params?: object) => [...notificationKeys.lists(), params] as const,
  details: () => [...notificationKeys.all, 'detail'] as const,
  detail: (id: number) => [...notificationKeys.details(), id] as const,
  unreadCount: () => [...notificationKeys.all, 'unreadCount'] as const,
  unread: () => [...notificationKeys.all, 'unread'] as const,
  preview: () => [...notificationKeys.all, 'preview'] as const,
}
