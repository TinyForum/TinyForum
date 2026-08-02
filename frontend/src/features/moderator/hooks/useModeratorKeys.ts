// hooks/useModeratorKeys.ts
export const moderatorKeys = {
  all: ["moderator"] as const,
  myApplications: (params?: object) =>
    [...moderatorKeys.all, "my-applications", params] as const,
  list: (boardId: number) => [...moderatorKeys.all, "list", boardId] as const,
  myBoards: (userId?: number) =>
    [...moderatorKeys.all, "my-boards", userId] as const,
  posts: (boardId: number, page?: number, keyword?: string) =>
    [...moderatorKeys.all, "posts", boardId, page, keyword] as const,
  reports: (boardId: number, page?: number) =>
    [...moderatorKeys.all, "reports", boardId, page] as const,
  bannedUsers: (boardId: number, page?: number) =>
    [...moderatorKeys.all, "banned-users", boardId, page] as const,
}
