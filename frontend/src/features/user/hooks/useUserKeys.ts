// hooks/useUserKeys.ts
export const userKeys = {
  all: ["users"] as const,
  lists: () => [...userKeys.all, "list"] as const,
  list: (params?: object) => [...userKeys.lists(), params] as const,
  details: () => [...userKeys.all, "detail"] as const,
  detail: (id: number) => [...userKeys.details(), id] as const,
  // 更细粒度的 key 工厂
  stats: () => [...userKeys.all, "stats"] as const,
  role: () => [...userKeys.all, "role"] as const,
  violations: () => [...userKeys.all, "violations"] as const,
  violationDetail: (id: string) => [...userKeys.all, "violation", id] as const,
  followers: (id: number, params?: object) =>
    [...userKeys.all, "followers", id, params] as const,
  following: (id: number, params?: object) =>
    [...userKeys.all, "following", id, params] as const,
  leaderboard: (simple: boolean, limit?: number) =>
    [...userKeys.all, "leaderboard", simple, limit] as const,
  followStatus: (userId: number, currentUserId?: number) =>
    [...userKeys.all, "followStatus", userId, currentUserId] as const,
}
