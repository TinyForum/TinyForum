// hooks/useAnnouncementKeys.ts
export const announcementKeys = {
  all: ['announcements'] as const,
  lists: () => [...announcementKeys.all, 'list'] as const,
  list: (params?: object) => [...announcementKeys.lists(), params] as const,
  details: () => [...announcementKeys.all, 'detail'] as const,
  detail: (id: number) => [...announcementKeys.details(), id] as const,
  pinned: (boardId?: number) => [
    ...announcementKeys.all,
    'pinned',
    boardId,
  ] as const,
}
