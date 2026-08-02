// hooks/useCommentKeys.ts
export const commentKeys = {
  all: ["comments"] as const,
  lists: () => [...commentKeys.all, "list"] as const,
  list: (postId: number, params?: { page?: number; page_size?: number }) =>
    [...commentKeys.lists(), postId, params] as const,
  details: () => [...commentKeys.all, "detail"] as const,
  detail: (id: number) => [...commentKeys.details(), id] as const,
};
